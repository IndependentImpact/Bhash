"""Create a Hedera Consensus Service topic and persist metadata in Fluree.

This helper wires together the Hedera SDK and the repository's lightweight
Fluree client so contributors can bootstrap a ledger that records consensus
topics.  The script expects Hedera operator credentials as well as Fluree
Cloud API secrets to be available in the environment and will abort with
helpful error messages when anything is missing.

Usage examples
--------------

Create a new topic, materialise a ledger called ``hedera-topics-<timestamp>``
under your Fluree tenant, and store the topic details::

    python scripts/hedera_topic_to_fluree.py

Reuse an existing ledger without attempting to create a dataset first::

    python scripts/hedera_topic_to_fluree.py --ledger my-handle/hedera-topics

Supply a custom dataset name and visibility while keeping the automatic
timestamp suffix::

    python scripts/hedera_topic_to_fluree.py --dataset-name hedera-topics \
        --visibility public

The script prints a short summary on success and will raise informative
exceptions whenever the Hedera or Fluree APIs return an error.  No attempt is
made to delete the topic afterwards – running the script multiple times will
therefore create multiple consensus topics.
"""

from __future__ import annotations

import argparse
import datetime as _dt
import logging
import os
import sys
from dataclasses import dataclass
from typing import Any, Dict, Iterable, Mapping, MutableMapping, Optional

if __package__ in (None, ""):
    # Allow running the script directly without installing the package.
    sys.path.append(os.path.dirname(os.path.dirname(__file__)))

from hedera import AccountId, Client, PrivateKey, TopicCreateTransaction

try:  # pragma: no cover - transitional compatibility layer
    from scripts.fluree_client import FlureeClient, FlureeClientError, FlureeConfig  # type: ignore
except ModuleNotFoundError as exc:  # pragma: no cover - transitional compatibility layer
    raise RuntimeError(
        "The Python Fluree client has been replaced by the Go implementation. "
        "Use 'go run ./cmd/bhashctl hedera bootstrap' or the scripts/flureeclient package instead."
    ) from exc

LOGGER = logging.getLogger("hedera_fluree")

# Compact JSON-LD context used when inserting topic metadata into Fluree.
JSONLD_CONTEXT: Mapping[str, Any] = {
    "@version": 1.1,
    "@vocab": "https://hashgraphontology.xyx/hedera/ledger#",
    "hedera": "https://hashgraphontology.xyx/hedera/core/",
    "xsd": "http://www.w3.org/2001/XMLSchema#",
    "Topic": "hedera:ConsensusTopic",
    "topicId": "https://hashgraphontology.xyx/hedera/ledger#topicId",
    "network": "hedera:occursOn",
    "memo": "https://hashgraphontology.xyx/hedera/ledger#memo",
    "consensusTimestamp": {
        "@id": "https://hashgraphontology.xyx/hedera/ledger#consensusTimestamp",
        "@type": "xsd:dateTime",
    },
    "transactionId": "https://hashgraphontology.xyx/hedera/ledger#transactionId",
}


@dataclass(frozen=True)
class TopicMetadata:
    """Normalised Hedera topic attributes suitable for JSON serialisation."""

    topic_id: str
    network: str
    consensus_timestamp: str
    transaction_id: str
    memo: Optional[str]

    def to_jsonld(self) -> MutableMapping[str, Any]:
        """Return a JSON-LD resource describing the topic."""

        identifier = self.topic_id.replace(".", "-")
        resource: MutableMapping[str, Any] = {
            "@id": f"topic/{identifier}",
            "@type": "Topic",
            "topicId": self.topic_id,
            "network": self.network,
            "consensusTimestamp": self.consensus_timestamp,
            "transactionId": self.transaction_id,
        }
        if self.memo:
            resource["memo"] = self.memo
        return resource


def _load_hedera_client() -> Client:
    """Initialise a Hedera client from environment variables."""

    network = os.getenv("HEDERA_NETWORK", "testnet").strip().lower()
    operator_id = os.getenv("HEDERA_OPERATOR_ID")
    operator_key = os.getenv("HEDERA_OPERATOR_KEY")

    missing = [
        name
        for name, value in {
            "HEDERA_OPERATOR_ID": operator_id,
            "HEDERA_OPERATOR_KEY": operator_key,
        }.items()
        if not value
    ]
    if missing:
        raise RuntimeError(
            "Missing Hedera credentials: " + ", ".join(sorted(missing))
        )

    try:
        if network == "mainnet":
            client = Client.forMainnet()
        elif network == "previewnet":
            client = Client.forPreviewnet()
        else:
            client = Client.forTestnet()
    except Exception as exc:  # pragma: no cover - defensive guard
        raise RuntimeError(f"Unsupported Hedera network '{network}': {exc}") from exc

    client.setOperator(AccountId.fromString(operator_id), PrivateKey.fromString(operator_key))
    return client


def _to_string(value: Any) -> str:
    """Convert Hedera SDK Java proxy objects into Python strings."""

    if value is None:
        return ""
    if hasattr(value, "toString"):
        try:
            return value.toString()
        except Exception:  # pragma: no cover - defensive guard
            return str(value)
    return str(value)


def _to_optional_string(value: Any) -> Optional[str]:
    text = _to_string(value)
    return text or None


def _network_name(client: Client) -> str:
    """Return the Hedera network name associated with the client."""

    return _to_string(client.getNetworkName()).lower()


def _build_topic_memo(custom_memo: Optional[str]) -> str:
    """Return a memo string that fits within the SDK limits."""

    base = custom_memo or "Bhash integration topic"
    encoded = base.encode("utf-8")
    if len(encoded) <= 100:
        return base
    LOGGER.warning("Memo too long (%d bytes), truncating", len(encoded))
    # The SDK enforces a 100 byte limit; truncate on the byte representation.
    return encoded[:100].decode("utf-8", errors="ignore")


def _create_topic(client: Client, memo: Optional[str]) -> TopicMetadata:
    """Create a consensus topic and capture essential metadata."""

    transaction = TopicCreateTransaction()
    if memo:
        transaction.setTopicMemo(memo)

    LOGGER.info("Submitting Hedera TopicCreateTransaction")
    response = transaction.execute(client)
    receipt = response.getReceipt(client)
    record = response.getRecord(client)

    topic_id = _to_string(receipt.topicId)
    consensus_timestamp = _to_string(record.consensusTimestamp)
    transaction_id = _to_string(record.transactionId)

    LOGGER.info("Created topic %s at %s", topic_id, consensus_timestamp)
    return TopicMetadata(
        topic_id=topic_id,
        network=_network_name(client),
        consensus_timestamp=consensus_timestamp,
        transaction_id=transaction_id,
        memo=_to_optional_string(getattr(record, "transactionMemo", None)) or memo,
    )


def _default_dataset_name(base_name: Optional[str]) -> str:
    """Return a dataset name suffixed with an ISO timestamp."""

    timestamp = _dt.datetime.now(tz=_dt.timezone.utc).strftime("%Y%m%d-%H%M%S")
    base = base_name or "hedera-topics"
    return f"{base}-{timestamp}"


def _ensure_ledger(
    client: FlureeClient,
    *,
    dataset_name: str,
    storage_type: str,
    description: str,
    visibility: str,
    tags: Optional[Iterable[str]] = None,
) -> str:
    """Create the Fluree dataset if it does not already exist."""

    owner = client.config.tenant_handle
    LOGGER.info("Ensuring dataset %s under handle %s", dataset_name, owner)
    try:
        client.create_dataset(
            owner,
            dataset_name=dataset_name,
            storage_type=storage_type,
            description=description,
            visibility=visibility,
            tags=tags,
        )
        LOGGER.info("Created dataset %s", dataset_name)
    except FlureeClientError as exc:
        message = str(exc).lower()
        if "already exists" in message:
            LOGGER.info("Dataset %s already exists, continuing", dataset_name)
        else:
            raise
    return f"{owner}/{dataset_name}"


def _store_topic_metadata(
    client: FlureeClient,
    *,
    ledger: str,
    metadata: TopicMetadata,
) -> Mapping[str, Any]:
    """Persist topic metadata to the Fluree ledger."""

    LOGGER.info("Transacting topic metadata into %s", ledger)
    insert = [metadata.to_jsonld()]
    return client.transact(ledger=ledger, insert=insert, context=JSONLD_CONTEXT)


def parse_arguments(argv: Optional[Iterable[str]] = None) -> argparse.Namespace:
    """Parse command-line arguments."""

    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--dataset-name",
        help="Base dataset name to create; a timestamp suffix is appended automatically.",
    )
    parser.add_argument(
        "--ledger",
        help=(
            "Fully qualified ledger identifier (handle/dataset). "
            "When provided, dataset creation is skipped."
        ),
    )
    parser.add_argument(
        "--memo",
        help="Custom topic memo (defaults to a descriptive placeholder).",
    )
    parser.add_argument(
        "--visibility",
        default="private",
        choices=["private", "public"],
        help="Fluree dataset visibility (default: private).",
    )
    parser.add_argument(
        "--storage-type",
        default=os.getenv("FLUREE_STORAGE_TYPE", "immutable"),
        help="Fluree storage type (defaults to the FLUREE_STORAGE_TYPE env var or 'immutable').",
    )
    parser.add_argument(
        "--tag",
        action="append",
        dest="tags",
        help="Optional tag to attach to the dataset (may be repeated).",
    )
    parser.add_argument(
        "--description",
        default="Ledger capturing Hedera consensus topics created by the Bhash toolkit.",
        help="Dataset description used when creating the ledger.",
    )
    return parser.parse_args(argv)


def main(argv: Optional[Iterable[str]] = None) -> int:
    """Entry point for the CLI."""

    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(message)s")
    args = parse_arguments(argv)

    memo = _build_topic_memo(args.memo)
    with _load_hedera_client() as hedera_client:
        metadata = _create_topic(hedera_client, memo)

    fluree_config = FlureeConfig.from_env()
    fluree_client = FlureeClient(fluree_config)

    if args.ledger:
        ledger = args.ledger
    else:
        dataset_name = _default_dataset_name(args.dataset_name)
        ledger = _ensure_ledger(
            fluree_client,
            dataset_name=dataset_name,
            storage_type=args.storage_type,
            description=args.description,
            visibility=args.visibility,
            tags=args.tags,
        )

    response = _store_topic_metadata(fluree_client, ledger=ledger, metadata=metadata)

    LOGGER.info("Stored topic metadata in %s", ledger)
    print(
        "Created Hedera topic {topic} on {network} and stored metadata in {ledger}.".format(
            topic=metadata.topic_id, network=metadata.network, ledger=ledger
        )
    )
    if response:
        print("Fluree response:")
        print(response)
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except FlureeClientError as exc:
        LOGGER.error("Fluree API error: %s", exc)
        sys.exit(2)
    except Exception as exc:  # pragma: no cover - user feedback
        LOGGER.error("%s", exc)
        sys.exit(1)

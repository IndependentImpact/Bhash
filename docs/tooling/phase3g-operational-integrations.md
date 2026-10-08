# Phase 3G Operational Integrations

Phase 3G connects the ontology-driven automation pipeline to live (or simulated) Hedera
operations and Fluree persistence flows. This runbook introduces the `bhashctl hedera`
command family, reusable bootstrap specifications, and JSON-LD payloads aligned with the
ontology so downstream datasets can be loaded without ad-hoc scripting.

## 1. Bootstrap workflow overview

```
$ go run ./cmd/bhashctl hedera bootstrap \
    --spec data/fluree/phase3g-bootstrap-spec.json \
    --ledger tenant/dataset-handle \
    --simulate
```

The command performs the following steps:

1. **Load bootstrap spec** – reads accounts, topics, and tokens from the JSON spec. The
   repository ships with `data/fluree/phase3g-bootstrap-spec.json` which demonstrates a
   treasury account, consensus topic, and fungible token aligned with the Phase 3 ontology.
2. **Provision artefacts** – uses the Hedera network abstraction (`internal/hedera`) to
   mint the requested artefacts. When `--simulate` is enabled (default) a deterministic
   mock network issues IDs and timestamps so fixtures can be regenerated locally.
3. **Assemble JSON-LD** – converts the resulting metadata into ontology-aligned JSON-LD
   nodes and composes a Fluree transaction payload ready for `/fluree/transact`.
4. **Persist to Fluree (optional)** – when `--commit` is supplied, the command authenticates
   using `FLUREE_API_TOKEN`/`FLUREE_HANDLE` (or flag overrides) and posts the transaction
   using the HTTP client under `internal/fluree`.

The command prints a structured JSON summary containing the network metadata, generated
artefacts, and the transaction payload so it can be inspected before submission.

## 2. Specification format

Bootstrap specifications capture the artefacts to be created and the ontology metadata
that should be exported. The structure matches the `internal/hedera.BootstrapSpec` type:

```json
{
  "network": "testnet",
  "ledger": "phase3g/hedera-bootstrap",
  "accounts": [
    {
      "alias": "phase3g-treasury",
      "memo": "Treasury account for Phase 3G fixtures",
      "publicKey": "...",
      "initialBalanceTinybar": 100000000,
      "tags": ["phase-3g", "bootstrap"]
    }
  ],
  "topics": [
    {
      "alias": "phase3g-telemetry",
      "memo": "Phase 3G consensus topic for telemetry mirroring",
      "tags": ["consensus", "telemetry"]
    }
  ],
  "tokens": [
    {
      "alias": "phase3g-token",
      "name": "Phase 3G Integration Token",
      "symbol": "P3G",
      "memo": "Demonstration token minted during the Phase 3G bootstrap",
      "treasuryAlias": "phase3g-treasury",
      "decimals": 2,
      "initialSupply": 100000,
      "supplyType": "INFINITE",
      "tokenType": "FUNGIBLE_COMMON",
      "tags": ["token-service", "bootstrap"]
    }
  ]
}
```

* `accounts` entries define treasury or operator accounts and capture memos/tags for the
  ontology.
* `topics` specify consensus topics used for telemetry or governance coordination.
* `tokens` describe fungible/non-fungible tokens and reference accounts via `treasuryAlias`.

## 3. Hedera network abstraction

The `internal/hedera` package introduces:

* **Config parsing** – `Config` reads `HEDERA_*` environment variables and supports
  overrides via CLI flags.
* **Network interface** – `Network` exposes `CreateAccount`, `CreateTopic`, and
  `CreateToken` with mock and SDK-backed implementations.
* **Bootstrapper** – orchestrates provisioning and resolves aliases so tokens reference
  their treasury accounts automatically.
* **JSON-LD export** – converts the resulting records into a Fluree transaction with a
  shared context (`hedera`, `ledger`, `prov`, `schema`, and `xsd` prefixes). Nodes use the
  canonical resource IRIs (`https://hashgraphontology.xyz/resource/{network}/{account|topic|token}/{shard}.{realm}.{num}`,
  D-0002) and link to the network's named individual (`hedera:Mainnet`, `hedera:Testnet`,
  or `hedera:Previewnet`, D-0005) with `hedera:registeredIn`; a network without a named
  individual, or an ID that is not `shard.realm.num`, is rejected.

  | Record field | Emitted as |
  | ------------ | ---------- |
  | account / topic / token ID | `hedera:hasAccountId` / `hedera:hasTopicId` / `hedera:hasTokenId` |
  | token type | `@type` `hedera:FungibleToken` (also for an empty type) or `hedera:NonFungibleToken`; an unrecognised type stays `hedera:Token` plus `ledger:tokenType` |
  | token symbol, decimals, treasury, max supply | `hedera:hasSymbol`, `hedera:hasDecimals`, `hedera:hasTreasury` (account IRI), `hedera:hasMaxSupply` |
  | token initial supply | `hedera:hasInitialSupply` (fungible tokens only) |
  | alias, token name, memo, tags | `schema:name`, `schema:name`, `schema:description`, `schema:keywords` |
  | creation time | `prov:generatedAtTime` |
  | account public key, topic initial sequence, token supply type | `ledger:publicKey`, `ledger:initialSequence`, `ledger:supplyType` |

  The `ledger:` namespace (`https://hashgraphontology.xyz/ledger#`, also used by
  `scripts/hedera_topic_to_fluree.py`) holds fields that have no ontology property yet; its
  terms are not part of the ontology.

When `--simulate=false`, the CLI instantiates the SDK-backed network and expects
`HEDERA_OPERATOR_ID`/`HEDERA_OPERATOR_KEY` to be present so real transactions can be
submitted. Mirror network URLs can be overridden via `--mirror-url`.

## 4. Fluree integration

The Fluree HTTP client (`internal/fluree`) is reused to submit the JSON-LD payload. Provide
credentials using environment variables or flags:

```
$ export FLUREE_API_TOKEN=...
$ export FLUREE_HANDLE=...
$ go run ./cmd/bhashctl hedera bootstrap \
      --spec data/fluree/phase3g-bootstrap-spec.json \
      --ledger tenant/dataset-handle \
      --commit
```

On success the command returns the Fluree response alongside the transaction payload,
allowing operators to confirm commit hashes or ledger identifiers.

## 5. Next steps

* Extend the bootstrap spec with additional artefacts (e.g., scheduled transactions or
  HIP-540 token parameters) as new ontology modules land.
* Generate derived fixtures for regression tests by persisting the JSON output under
  `tests/fixtures/` and exercising the data through SPARQL/SHACL pipelines.
* Integrate the command into CI smoke tests once credentials are available so Phase 3G
  runs can continuously verify the Hedera → Fluree pipeline.

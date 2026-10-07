# Core Competency Question Answers

Phase 2 establishes executable evidence for foundational reasoning tasks backed by the core ontology module.

## CQ-CORE-001 – Which entities participate in consensus on mainnet?

* **Status:** Example data and SPARQL query committed.
* **Scope:** Enumerate mainnet validators, their controlling accounts, and the governance mandate that authorised their
  onboarding.
* **Inputs:**
  * Ontology: `ontology/src/core.ttl`
  * Sample graph: `ontology/examples/core-consensus.ttl`
  * Query: `tests/queries/cq-core-001.rq`

### Execution notes

1. Load the core ontology and example graph into your RDF store or run a temporary query using Apache Jena ARQ:
   ```bash
   arq --data ontology/src/core.ttl --data ontology/examples/core-consensus.ttl --query tests/queries/cq-core-001.rq
   ```
2. The expected bindings return the Hedera mainnet, consensus node 3, its staking account, and the council mandate citing HIP-840.
3. Extend the dataset with live mirror-node exports to validate against production identifiers once automation tasks from
   `docs/tooling/automation-requirements.md` are implemented.

### Sample result (derived from `ontology/examples/core-consensus.ttl`)

| network | validator | operatorAccount | mandate |
| ------- | --------- | ---------------- | ------- |
| `hedera:Mainnet` | `ex:Node3` | `ex:Node3Account` | "HIP-840 validator onboarding" |

### Evidence bundle

Detailed documentation anchors, data sourcing, and modelling notes are captured in
[`docs/competency/evidence/CQ-CORE-001.md`](evidence/CQ-CORE-001.md).

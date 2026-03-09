# Modelling Decision Log

This log captures notable decisions, trade-offs, and open questions encountered while modelling the Hedera/Hiero ontology.  Each entry should link to the supporting documentation (e.g., bibliography citations, competency questions, pull requests) and record the status for traceability.

| ID | Date | Decision | Context & Rationale | Status | Owner |
| -- | ---- | -------- | ------------------- | ------ | ----- |
| D-0001 | 2024-05-01 | Adopt ROBOT for ontology automation | Compared ROBOT and RDFlib for build/test automation. ROBOT provides purpose-built OWL workflows (templates, reasoning, report generation) and integrates with CI pipelines for validation. RDFlib remains available for data scripting, but ROBOT will anchor automated builds. | Accepted | Ontology Engineering Team |
| D-0002 | 2026-03-09 | Adopt canonical instance IRI patterns for Hedera resources | Projects consuming the ontology were minting local IRIs for accounts, topics, messages, tokens, and contracts, hurting interoperability and making `owl:sameAs` alignment manual. A deterministic scheme under `https://hashgraphontology.xyz/resource/{network}/{type}/{shard}.{realm}.{num}` was adopted, anchored on native Hedera entity IDs. Topic messages use `…/message/{sequenceNumber}` as the stable sub-key. The patterns are encoded as `hedera:instanceIRIPattern` annotations on each class, enforced with SHACL `sh:Warning` shapes in `ontology/shapes/instance-iris.shacl.ttl`, and illustrated in `ontology/examples/canonical-iris.ttl`. See also the root `README.md` canonical IRI section. | Accepted | Ontology Engineering Team |

## How to propose updates

1. **Draft entry** – add a new row with a unique ID (`D-XXXX`), provisional status (`Proposed`), and reference links.
2. **Review** – discuss in issues or pull requests; capture alternatives considered.
3. **Finalise** – once approved, update the status to `Accepted` or `Rejected` and reference the decision artefact.

Maintain chronological order and ensure each decision points back to the relevant sources in `docs/references/bibliography.md`.

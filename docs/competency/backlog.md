# Phase 1 Competency Question Backlog

This backlog captures high-priority competency questions (CQs) that will guide ontology modelling, SHACL validation, and SPARQL query development.  Each entry references the stakeholder motivation, expected data sources, and relevant documentation anchors.

## How to use this backlog

1. Reference IDs in issues, ontology modules, and SHACL shapes so traceability is preserved.
2. Update the **Status** column as questions move from draft to validated (via SPARQL or SHACL evidence).
3. Link answers or query artefacts once produced to accelerate regression testing.

| ID | Stakeholder | Priority | Question | Source(s) | Notes | Status |
| -- | ----------- | -------- | -------- | --------- | ----- | ------ |
| CQ-GOV-001 | Governance | High | Which Hedera Council members currently steward validator onboarding decisions and what HIPs underpin their mandates? | [Hedera Council](https://hedera.com/council); [HIP-840](https://hips.hedera.com/hip/hip-840) | Requires modelling council committees, meeting records, and HIP provenance; evidence bundle captured in [`evidence/CQ-GOV-001.md`](evidence/CQ-GOV-001.md). | Validated |
| CQ-GOV-002 | Governance | High | What quorum of council votes authorised the latest network fee schedule update? | [Council governance](https://hedera.com/council); [Fee schedule docs](https://docs.hedera.com/hedera/core-concepts/fees) | Dependent on availability of public vote records or release notes. | Draft |
| CQ-COMP-003 | Compliance | High | Which tokens classified as stablecoins (HIP-540) enforce KYC and freeze keys, and who controls those keys? | [HIP-540](https://hips.hedera.com/hip/hip-540); [Token Service](https://docs.hedera.com/hedera/sdks-and-apis/token-service/introduction) | Drives modelling of token compliance attributes and account roles; evidence bundle stored in [`evidence/CQ-COMP-003.md`](evidence/CQ-COMP-003.md). | In review |
| CQ-COMP-004 | Compliance | Medium | Which scheduled transactions remain pending beyond 24 hours due to missing signatures? | [Scheduled transactions](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions); [Mirror REST schedules](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Schedule) | SHACL and SPARQL assets committed; awaiting live mirror data feeds. | In review |
| CQ-DEV-005 | Developer tooling | High | Which smart contracts invoke HTS system contract precompiles and what gas usage patterns do they exhibit? | [HSCS system contracts](https://docs.hedera.com/hedera/core-concepts/smart-contracts/system-contracts); [Mirror contract logs](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Contracts) | Supports optimisation guidance and ontology alignment with execution traces; evidence bundle recorded in [`evidence/CQ-DEV-005.md`](evidence/CQ-DEV-005.md). | In review |
| CQ-DEV-006 | Developer tooling | Medium | Which topics enforce both admin and submit key signatures for configuration updates? | [HCS configuration](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service/manage-topics); [Mirror topics endpoint](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Topics) | Connects governance metadata with operational artefacts. | Draft |
| CQ-ANL-007 | Analytics | High | How is total supply of fungible tokens distributed across treasury and external accounts over time? | [Token Service](https://docs.hedera.com/hedera/sdks-and-apis/token-service/token-relationships); [Mirror balances](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Balances) | Dataset retention modelled; quantitative aggregation pending real balance extracts. | In review |
| CQ-ANL-008 | Analytics | Medium | Which consensus topics have retention periods shorter than 30 days and which dApps depend on them? | [HCS retention](https://docs.hedera.com/hedera/core-concepts/hedera-consensus-service/message-retention); [Mirror subscriptions](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/topics) | Encourages modelling of topic consumers and SLAs. | Draft |
| CQ-IMPACT-001 | Sustainability | High | How do scheduled sustainability transactions translate into AIAO impact claims with supporting KPIs and telemetry provenance? | [Scheduled transactions](https://docs.hedera.com/hedera/core-concepts/scheduled-transactions); [Mirror metrics](https://docs.hedera.com/hedera/mirror-node/architecture/metrics); [AIAO](https://w3id.org/aiao) | Alignment query and example graph validated in Phase 4; bridge modules published for anthropogenic impact onboarding. | Validated |
| CQ-HIE-009 | Hiero transition | High | Which validators are participating in each Hiero shard and what onboarding state (pending, active, slashed) do they occupy? | [HIP-840](https://hips.hedera.com/hip/hip-840); [Hiero docs](https://docs.hedera.com/hiero) | Ontology, SHACL, and query assets landed; awaiting Hiero registry feeds. | In review |
| CQ-HIE-010 | Hiero transition | Medium | How do shard-level consensus finality times compare to mainnet benchmarks for a given epoch? | [Hiero docs](https://docs.hedera.com/hiero); [Mirror performance metrics](https://docs.hedera.com/hedera/mirror-node/architecture/metrics) | Will require time-series data capture and modelling of performance indicators. | Draft |

## High-priority execution plan

The following work packages decompose the three "High" priority competency questions into actionable issues. Each task should be
captured as a GitHub issue (label: `competency`) with acceptance criteria reflecting the deliverables below.

| Task ID | CQ | Focus | Deliverables | Dependencies | Status |
| ------- | -- | ----- | ------------ | ------------ | ------ |
| CQ-GOV-001-A | CQ-GOV-001 | Council roster ingestion | `data/council-roster.csv` sourced from hedera.com with committee assignments, plus provenance notes in [`evidence/CQ-GOV-001.md`](evidence/CQ-GOV-001.md). | Access to council roster; scraping script. | ✅ Complete |
| CQ-GOV-001-B | CQ-GOV-001 | Governance modelling | Update consensus module with validator stewardship roles and committee hierarchies aligned to glossary entries. | CQ-GOV-001-A | ✅ Complete |
| CQ-GOV-001-C | CQ-GOV-001 | Validation assets | SPARQL query committed to `tests/queries/cq-gov-001.rq` plus SHACL shape verifying node–steward linkage. | CQ-GOV-001-B; automation tasks AUT-004/005. | ✅ Complete |
| CQ-COMP-003-A | CQ-COMP-003 | Token compliance snapshot | Mirror REST extraction script for stablecoin tokens with key metadata stored in `data/token-compliance.json`. | Mirror node API access. | In review |
| CQ-COMP-003-B | CQ-COMP-003 | Role alignment | Extend glossary and ontology annotations to map token admin/freeze/metadata custodians to `hedera:ComplianceSteward` subclasses. | CQ-COMP-003-A; glossary updates. | Not started |
| CQ-COMP-003-C | CQ-COMP-003 | Query & SHACL | SPARQL query `tests/queries/cq-comp-003.rq` and SHACL rule ensuring every stablecoin exposes accountable key custodians. | CQ-COMP-003-B; AUT-004/005. | ✅ Complete |
| CQ-DEV-005-A | CQ-DEV-005 | Contract log fixture | Download sample contract log dataset (JSON/CSV) into `data/contracts/hts-precompiles/` with README. | Mirror node contract log endpoint; storage location. | 🔄 In progress |
| CQ-DEV-005-B | CQ-DEV-005 | Ontology extension | Introduce `hedera:PrecompileInvocation` class and link contracts to HTS precompile terms in `ontology/src/smart-contracts.ttl`. | CQ-DEV-005-A; core ontology module. | ✅ Complete |
| CQ-DEV-005-C | CQ-DEV-005 | Analytics validation | SPARQL aggregation `tests/queries/cq-dev-005.rq` summarising gas usage and invocation counts, with expected results stored in `tests/fixtures/cq-dev-005.json`. | CQ-DEV-005-B; AUT-005. | In review |

## Next actions

* Prioritise CQ-GOV-001, CQ-COMP-003, and CQ-DEV-005 for initial ontology modelling, as they cover governance, compliance, and developer tooling focal areas.
* Identify authoritative data sources (council minutes, HIP releases, mirror node datasets) to ground each question before drafting SPARQL queries.
* Add lower-priority questions (e.g., education, ecosystem tooling) once foundational modules stabilise.

# External Ontology Landscape Review (Phase 1)

Phase 1 requires identifying reusable vocabularies that can accelerate Bhash modelling.  This review surveys candidate ontologies and patterns, highlighting immediate reuse opportunities and follow-up actions.

## Summary matrix

| Ontology / Pattern | Domain coverage | Reuse opportunities | Considerations |
| ------------------ | --------------- | ------------------- | --------------- |
| [PROV-O](https://www.w3.org/TR/prov-o/) | Provenance of activities, agents, and entities. | Transaction lifecycle, governance decisions, mirror export derivations. | Determine profiling for consensus vs execution activities. |
| [DCAT 3](https://www.w3.org/TR/vocab-dcat-3/) | Dataset cataloguing. | Mirror node REST/gRPC datasets, record/balance file descriptions. | Extend with Hedera-specific distribution metadata (e.g., network, retention). |
| [DID Core](https://www.w3.org/TR/did-core/) | Decentralised identifier representation. | Accounts, keys, and wallet bindings. | Assess compatibility with Hedera account ID format and key rotation semantics. |
| [FIBO](https://spec.edmcouncil.org/fibo/) | Financial instruments & governance. | Token classifications, treasury roles, governance committees. | Scope alignment needed; avoid over-constraining token semantics. |
| [ODRL](https://www.w3.org/TR/odrl-model/) | Policy/permission modelling. | Custom fee rules, access control statements for topics/files. | Evaluate complexity relative to simpler policy shapes. |
| [OpenTelemetry semantic conventions](https://opentelemetry.io/docs/specs/semconv/) | Telemetry schema for observability. | Potential mapping for node metrics and mirror observability data. | Likely optional until telemetry integration is prioritised. |
| [AIAO](https://w3id.org/aiao) | Anthropogenic impact accounting upper ontology. | Provides environmental and social impact scaffolding for anthropogenic impact reporting on Hedera-hosted assets. | Requires bridge module translating Hedera service events into impact indicators. |
| [ClaimOnt](https://datadudes.xyz/claimont) | Climate accounting ontology for mitigation/adaptation measures. | Align Hedera staking, sustainability initiatives, and carbon credit tokens with climate metrics. | Need domain expert review to avoid overstating equivalence relationships. |
| [ImpactOnt](https://datadudes.xyz/impactont) | Ontology for modelling impact investment portfolios and KPIs. | Map Hedera token compliance and treasury analytics concepts to impact investment observables. | Ensure financial definitions complement existing FIBO alignments. |
| [InfoComm](https://datadudes.xyz/infocomm) | Information and communication infrastructure ontology. | Useful for representing Hedera network infrastructure, mirror data pipelines, and observability dependencies. | Identify overlap with Hiero virtualization constructs before importing. |

## Detailed notes

### PROV-O

* Map Hedera transactions to `prov:Activity` instances with associated `prov:Agent` roles (payer, operator, validator).
* Represent mirror node exports as `prov:Entity` derived from record files, enabling traceability for analytics pipelines.
* Next step: design a minimal PROV profile capturing transaction submission, consensus finalisation, execution, and export.

### DCAT 3

* Use `dcat:Dataset` for mirror REST/gRPC datasets (transactions, balances, topics) and `dcat:Distribution` for endpoint-specific metadata.
* Capture shard or network scope via `dcat:spatial` or custom properties referencing Hiero shards.
* Next step: draft a template for dataset metadata referencing the [Mirror node REST API](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api).

### DID Core

* Hedera accounts and keys can be modelled as DID subjects with verification methods referencing ED25519/threshold key material.
* Wallet integrations can express account custody and key rotation via DID Document updates.
* Next step: evaluate whether Hedera's `0.0.x` account identifiers should appear as DID methods (e.g., `did:hedera:`) or remain separate identifiers.

### FIBO

* Provides rich financial governance and organisational vocabulary for council committees, token categories (stablecoin, asset-backed), and compliance roles.
* Complex ontology; consider importing only relevant modules (e.g., `FND`, `IND`, `BUs`) to avoid reasoning overhead.
* Next step: create a mapping table for token and governance terms against FIBO classes referenced in the [Token Service documentation](https://docs.hedera.com/hedera/sdks-and-apis/token-service/introduction).

### ODRL

* Offers policy expressions that can model token custom fees (constraints, duties, permissions) and file/topic access control statements.
* Might be heavyweight for early phases; consider SHACL-based policy shapes first.
* Next step: prototype an ODRL policy capturing an HTS royalty rule defined in [HIP-423](https://hips.hedera.com/hip/hip-423).

### OpenTelemetry semantic conventions

* Relevant for future monitoring use cases (node metrics, mirror API performance) mentioned in [Hedera mirror node metrics docs](https://docs.hedera.com/hedera/mirror-node/architecture/metrics).
* Could align Hiero observability artefacts with standard telemetry terms.
* Next step: catalogue which metrics are publicly exposed to gauge scope.

### AIAO

* Emphasises anthropogenic impact indicators (environmental, social, governance) that downstream reporting platforms require.
* Hedera token and treasury analytics modules can expose activity traces (e.g., sustainability-linked tokens, carbon offsets) as `aiao:ImpactClaim` instances referencing consensus transactions.
* Next step: design a bridge shape that consumes Hedera transaction provenance and emits AIAO-compatible impact statements for anthropogenic impact pilots.

### ClaimOnt

* Focused on climate mitigation/adaptation measures, emissions accounting, and policy instruments.
* Scheduled transactions and governance artefacts can document climate-related commitments by mapping to `claimont:ClimateAction` and `claimont:MitigationMeasure` classes.
* Next step: convene with Hedera sustainability programme owners to validate which ledger events constitute authoritative climate actions.

### ImpactOnt

* Models impact investment theses, KPIs, and measurement frameworks used by anthropogenic impact funds.
* HTS compliance and treasury analytics data can populate `impactont:ImpactMetric` observations, linking to underlying token supply/movement events.
* Next step: create competency questions covering impact reporting (e.g., "Which Hedera-issued tokens contribute to SDG-aligned KPIs?") and prototype SPARQL mappings to ImpactOnt terms.

### InfoComm

* Provides vocabulary for communication infrastructure, data flows, and service dependencies.
* Aligns with Hedera mirror node and Hiero virtualization constructs, enabling explicit modelling of data pipelines and observability hooks.
* Next step: assess overlap with existing network topology classes in the core and mirror modules, then define equivalence or subClassOf relations where semantics align.

## Action items

1. Draft minimal import profiles for PROV-O and DCAT to include in the core ontology module.
2. Schedule deeper evaluation sessions for DID Core and FIBO once account/governance modelling begins.
3. Track experimental ODRL/telemetry alignments as optional extensions until concrete competency questions demand them.

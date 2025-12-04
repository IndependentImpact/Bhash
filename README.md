# Bhash: Hedera Network Ontology

Bhash is an open knowledge engineering effort to describe the Hedera Network and the forthcoming Hiero architecture in a machine-interpretable way. The project captures the actors, assets, services, processes, and governance rules that appear across:

* Hedera Consensus Service (HCS),
* Token Service (HTS),
* Smart Contract Service (HSCS),
* File Service,
* Scheduled Transactions,
* mirror node ecosystem, and
* Hiero overlay.

By modelling Hedera's public documentation and implementation guidance—including Hiero validator onboarding—we provide reusable semantics for analytics, compliance, education, and integration projects. Phase 4 extends this scope with bridge modules that align Hedera-native concepts with the 
* Anthropogenic Impact Accounting Ontology (AIAO),
* Claim Ontology,
* Impact Ontology, and
* Information Communication Ontology

## Why an ontology?

* **Shared vocabulary** – establish stable identifiers and definitions for Hedera-specific notions such as accounts, topics, tokens, scheduled transactions, staking nodes, and fee schedules.
* **Interoperability** – align Hedera concepts with standard vocabularies (PROV-O, DCAT, W3C DID Core, etc.) and specialised anthropogenic impact ontologies (AIAO, ClaimOnt, ImpactOnt, InfoComm) so that data from mirror nodes, dApps, sustainability disclosures, and compliance tools can interoperate.
* **Reasoning & validation** – enable automated validation (via OWL reasoning and SHACL constraints) for network states, policy rules, token compliance requirements, smart contract metadata, and Hiero onboarding milestones.
* **Documentation** – provide an authoritative reference that augments Hedera/Hiero manuals with explicit relationships that are otherwise scattered across prose and code.

## Current Phase 3 deliverables

Phase 3 targets service-specific ontology modules. The following artefacts are now available:

| Domain | Ontology module | SHACL | Example graph | Competency assets |
| ------ | ---------------- | ----- | ------------- | ----------------- |
| Consensus Service | `ontology/src/consensus.ttl` | `ontology/shapes/consensus.shacl.ttl` | `ontology/examples/core-consensus.ttl` | `docs/competency/core-foundational.md` |
| Token Service | `ontology/src/token.ttl` | `ontology/shapes/token.shacl.ttl` | `ontology/examples/token-compliance.ttl` | `docs/competency/token-compliance.md` |
| Smart Contract Service | `ontology/src/smart-contracts.ttl` | `ontology/shapes/smart-contracts.shacl.ttl` | `ontology/examples/smart-contracts.ttl` | `docs/competency/smart-contracts.md` |
| File & Schedule Services | `ontology/src/file-schedule.ttl` | `ontology/shapes/file-schedule.shacl.ttl` | `ontology/examples/file-schedule.ttl` | `docs/competency/file-schedule.md` |
| Mirror & Analytics Ecosystem | `ontology/src/mirror-analytics.ttl` | `ontology/shapes/mirror-analytics.shacl.ttl` | `ontology/examples/mirror-analytics.ttl` | `docs/competency/mirror-analytics.md` |
| Hiero overlay | `ontology/src/hiero.ttl` | `ontology/shapes/hiero.shacl.ttl` | `ontology/examples/hiero.ttl` | `docs/competency/hiero.md` |

Each module is backed by competency questions, SPARQL regression queries under `tests/queries/`, expected results in `tests/fixtures/results/`, and SHACL shapes verifying structural requirements.

## External ontology alignment

The Phase 4 integration sprint introduces dedicated bridge modules that expose Hedera artefacts to anthropogenic impact and communications teams via four priority ontologies:

| Bridge | File | Focus |
| ------ | ---- | ----- |
| AIAO alignment | `ontology/src/alignment/aiao.ttl` | Maps consensus evidence, scheduled transactions, and token reserve events to `aiao:ImpactClaim` patterns for sustainability attestations. |
| ClaimOnt alignment | `ontology/src/alignment/claimont.ttl` | Projects scheduled climate commitments and milestones into ClaimOnt mitigation/adaptation taxonomies. |
| ImpactOnt alignment | `ontology/src/alignment/impactont.ttl` | Aligns treasury and compliance telemetry with ImpactOnt KPI, SDG, and policy structures. |
| InfoComm alignment | `ontology/src/alignment/infocomm.ttl` | Describes mirror nodes, Hiero shards, and telemetry exchanges as InfoComm communication assets with latency/SLA annotations. |

All four modules reuse the same prefixes (`ontology/src/alignment/prefixes.ttl`), cite Hedera/Hiero documentation alongside external ontology references, and are exercised by the anthropogenic impact competency query (`tests/queries/cq-impact-001.rq`) plus the alignment example graph (`ontology/examples/alignment-impact.ttl`).

## Repository layout

```text
.
├── README.md
├── data/                     # Source datasets and fixtures backing competency questions
│   ├── contracts/hts-precompiles/sample-invocations.csv
│   ├── mirror/token-balance-retention.csv
│   └── token-compliance.json
├── docs/
│   ├── competency/           # Competency question answers, backlog, and evidence bundles
│   ├── mappings/             # Crosswalks linking ontology terms to source documentation
│   ├── workplan.md           # Iterative roadmap and phase reviews
│   └── ...                   # Research notes, references, governance decisions
├── ontology/
│   ├── src/                  # OWL/Turtle source files for each module
│   ├── shapes/               # SHACL shapes aligned with module requirements
│   └── examples/             # Example graphs powering SPARQL and SHACL regression tests
├── scripts/                  # Automation for running SPARQL and SHACL checks
└── tests/
    ├── queries/              # SPARQL regression queries
    ├── fixtures/results/     # Expected query outputs
    └── fixtures/datasets/    # Supplementary RDF datasets
```

## Running validation

Use the Go-based CLI to install validation tooling, interact with Fluree datasets, and
execute regression checks. The commands below download the required binaries on demand
and materialise results under `build/` (created automatically):

```bash
go run ./cmd/bhashctl install          # Fetch ROBOT + TopBraid SHACL into build/tools
go run ./cmd/bhashctl sparql           # Execute SPARQL regression queries via ROBOT
go run ./cmd/bhashctl shacl            # Run SHACL validation with the TopBraid CLI
go run ./cmd/bhashctl fluree transact  # Apply JSON-LD transactions to a Fluree ledger
go run ./cmd/bhashctl hedera bootstrap # Create Hedera artefacts and export ontology-aligned JSON-LD
```

The CLI reuses the repository fixtures and reports mismatches against expected
snapshots. A handful of legacy Python harnesses remain under `scripts/` while
their Go replacements are scheduled, but new automation should target the Go
workflow.

## Generating OWL/HTML/JSON-LD artefacts

Publishable artefacts for each ontology module (OWL/RDF/XML, HTML catalogue,
JSON-LD, and Turtle) are generated from the source `.ttl` files under
`ontology/src` via a small Python helper. The script renders human-readable
catalogues using Jinja2 and keeps machine-readable formats in sync for
deployments.

1. Prepare the Python environment using the helper script under `ontology/` (this creates `ontology/venv` with the required
   dependencies):

   ```bash
   ontology/install_requirements.sh
   ```

2. Convert the source ontologies into deployable artefacts:

   ```bash
   ontology/venv/bin/python ontology/scripts/convert_ontologies.py \
     --source-dir ontology/src \
     --deployment-dir ontology/deployment
   ```

   Each module gains four files in `ontology/deployment/`: `.owl` (RDF/XML),
   `.html` (rendered catalogue), `.jsonld` (JSON-LD context/graph), and `.ttl`
   (normalised Turtle). A shared `imports/` directory is also refreshed to keep
   imported slices aligned.

3. Deploy by syncing the `ontology/deployment/` directory to your hosting
   target (e.g., GitHub Pages, S3 bucket, or an internal web server). For a
   static publish you can simply run:

   ```bash
   rsync -avh ontology/deployment/ <user>@<host>:/var/www/ontology/
   ```

   Downstream systems that consume the OWL or JSON-LD contexts should reference
   the files from this deployment directory so they stay aligned with the
   repository sources.

## Phase 4 data pilot

Run `python scripts/run_phase4_pilot.py` to initialise an Oxigraph-backed triple store with the Phase 3/4 ontologies and example graphs. The script executes the anthropogenic impact competency query, runs SHACL validation, and records artefacts under `build/pilots/phase4/` for stakeholder review. A Go-native pilot command is on the roadmap; see `docs/python_helper_migration_plan.md` for the migration schedule.

## Supporting datasets

New fixtures unlock automation reuse across modules:

* `data/contracts/hts-precompiles/sample-invocations.csv` – representative mirror node export for HTS precompile analytics.
* `data/token-compliance.json` – snapshot of key custodians backing the HIP-540 example.
* `data/mirror/token-balance-retention.csv` – retention planning worksheet for treasury analytics datasets.

These lightweight fixtures inform the example RDF graphs and will be replaced with live extracts once data pipelines are automated.

## Working practices

1. **Document-first research** – extract canonical definitions from Hedera/Hiero documentation, HIPs, and mirror node references before introducing new classes.
2. **Iterative modelling** – deliver scoped ontology modules per Hedera service, validated with sample graphs and SPARQL competency queries.
3. **Community alignment** – involve Hedera developer relations, HIP authors, and compliance experts for terminology approval and governance modelling.
4. **Automation** – rely on the Go CLI (`go run ./cmd/bhashctl …`) to orchestrate ROBOT- and TopBraid-backed validation; legacy Python scripts remain only for archival reference.
5. **Versioning** – use semantic versioning for ontology releases with changelogs capturing class/property additions and deprecations.

## Getting involved

Contributions are welcome! Please open an issue describing the Hedera/Hiero concept you plan to model, include citations to relevant documentation, and propose competency questions your addition should satisfy. Pull requests should include updated OWL/Turtle files, documentation, mapping tables, and validation evidence (SPARQL outputs or SHACL reports).

## Reference materials

* [Hedera documentation](https://docs.hedera.com/hedera) – canonical service specifications, API references, and governance processes.
* [Hiero architecture overview](https://docs.hedera.com/hiero) – modular network evolution, shard design, and permissionless validator roadmap.
* [Hedera Improvement Proposals (HIPs)](https://hips.hedera.com/) – normative standards for tokens, accounts, mirror node APIs, and upcoming Hiero capabilities.

## License

This repository is distributed under the terms of the [Apache License 2.0](LICENSE).

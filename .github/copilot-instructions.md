# Copilot Instructions for Bhash

## Project Overview

Bhash is an open knowledge engineering project that models the Hedera Network and the Hiero architecture using machine-interpretable OWL ontologies. It provides:

- Semantic definitions for Hedera's services: HCS, HTS, HSCS, File Service, Scheduled Transactions, mirror node ecosystem, and Hiero overlay.
- SHACL validation shapes for network state and policy compliance.
- Phase 4 bridge modules aligning Hedera-native concepts with external ontologies (AIAO, ClaimOnt, ImpactOnt, InfoComm).
- A Go CLI tool (`bhashctl`) that orchestrates ROBOT- and TopBraid-backed validation, Fluree triple store transactions, and Hedera network interactions.

## Repository Layout

```
.
├── cmd/bhashctl/          # Go CLI entry points (main.go, hedera_cmd.go)
├── data/                  # Source CSV/JSON fixtures for competency questions
├── docs/
│   ├── competency/        # Competency question answers and evidence
│   ├── mappings/          # Term-to-documentation crosswalks
│   └── workplan.md        # Iterative roadmap and phase reviews
├── internal/
│   ├── fluree/            # Fluree triple store HTTP client
│   ├── hedera/            # Hedera SDK integration (sdk.go, config.go)
│   └── tools/             # ROBOT/SHACL/SPARQL orchestration wrappers
├── ontology/
│   ├── src/               # OWL/Turtle source modules (*.ttl)
│   ├── shapes/            # SHACL shapes (*.shacl.ttl)
│   ├── examples/          # Example RDF graphs
│   ├── deployment/        # Generated artefacts (OWL, HTML, JSON-LD, Turtle)
│   └── scripts/           # Python helpers for artefact generation
├── scripts/               # Python automation (run_shacl.py, run_sparql.py, etc.)
├── templates/             # ROBOT CSV templates
└── tests/
    ├── queries/           # SPARQL regression queries (*.rq)
    └── fixtures/
        ├── results/       # Expected SPARQL outputs
        └── datasets/      # RDF test datasets
```

## Tech Stack

- **Go 1.22** – CLI (`cmd/bhashctl/`) and tool orchestration (`internal/`). Module path: `github.com/hashgraph/bhash`.
- **OWL/Turtle** – Ontology source files in `ontology/src/`.
- **SHACL** – Validation shapes in `ontology/shapes/`.
- **SPARQL** – Regression queries in `tests/queries/`.
- **Python 3** – Legacy scripts under `scripts/` and `ontology/scripts/` (being migrated to Go).
- **ROBOT** – OWL reasoning and SPARQL execution (downloaded on demand by `bhashctl install`).
- **TopBraid SHACL CLI** – SHACL validation (downloaded on demand by `bhashctl install`).

## Key Commands

### Go CLI (primary workflow)

```bash
go run ./cmd/bhashctl install          # Download ROBOT + TopBraid SHACL into build/tools/
go run ./cmd/bhashctl sparql           # Run SPARQL regression queries via ROBOT
go run ./cmd/bhashctl shacl            # Run SHACL validation via TopBraid CLI
go run ./cmd/bhashctl fluree transact  # Apply JSON-LD transactions to a Fluree ledger
go run ./cmd/bhashctl hedera bootstrap # Create Hedera artefacts and export JSON-LD
```

### Go tests

```bash
go test ./...                          # Run all Go unit tests
go test ./internal/fluree/...          # Fluree client tests only
```

### Make targets

```bash
make all          # reason-core + report-core + shacl + sparql
make reason-core  # OWL reasoning with ELK reasoner
make report-core  # Generate ROBOT quality report
make shacl        # Run SHACL validation (uses Python venv)
make sparql       # Run SPARQL regression queries (uses Python venv)
make fluree-smoke # go test ./internal/fluree ./scripts/flureeclient
make clean        # Remove build/ directory
```

### Ontology artefact generation (Python)

```bash
ontology/install_requirements.sh      # Create ontology/venv with required packages
ontology/venv/bin/python ontology/scripts/convert_ontologies.py \
  --source-dir ontology/src \
  --deployment-dir ontology/deployment
```

## Coding Conventions

### Go

- Target **Go 1.22**; avoid features from later versions.
- Follow standard Go idioms (`gofmt`, `go vet`).
- New automation should target the Go CLI workflow; do **not** add new Python scripts unless explicitly asked.
- Place new CLI subcommands in `cmd/bhashctl/` and supporting logic in the appropriate `internal/` package.
- Use `github.com/rs/zerolog` for structured logging (already imported via the Hedera SDK).
- Error handling: return errors up the call stack with context using `fmt.Errorf("…: %w", err)`.

### Ontology / Turtle

- Use the existing prefix declarations in each module; do not introduce new namespaces without updating `ontology/src/alignment/prefixes.ttl` and all affected modules.
- Every new class or property must have `rdfs:label`, `rdfs:comment`, and `rdfs:isDefinedBy` annotations.
- Canonical instance IRIs follow the pattern `https://hashgraphontology.xyz/resource/{network}/{type}/{shard}.{realm}.{num}`.
- New modules must be accompanied by a SHACL shapes file in `ontology/shapes/` and at least one example graph in `ontology/examples/`.

### SPARQL / Tests

- Add competency questions as `.rq` files under `tests/queries/` with corresponding expected output snapshots in `tests/fixtures/results/`.
- Query file names follow the pattern `cq-<module>-<NNN>.rq`.

## Canonical IRI Patterns

| Resource | IRI pattern |
|----------|-------------|
| Account | `https://hashgraphontology.xyz/resource/{network}/account/{shard}.{realm}.{num}` |
| Consensus Topic | `https://hashgraphontology.xyz/resource/{network}/topic/{shard}.{realm}.{num}` |
| Topic Message | `https://hashgraphontology.xyz/resource/{network}/topic/{shard}.{realm}.{num}/message/{sequenceNumber}` |
| Token | `https://hashgraphontology.xyz/resource/{network}/token/{shard}.{realm}.{num}` |
| Smart Contract | `https://hashgraphontology.xyz/resource/{network}/contract/{shard}.{realm}.{num}` |

`{network}` is one of `mainnet`, `testnet`, or `previewnet`.

## Working Practices

1. **Document-first** – extract canonical definitions from Hedera/Hiero documentation, HIPs, and mirror node references before introducing new classes.
2. **Iterative modelling** – deliver scoped ontology modules per Hedera service, validated with sample graphs and SPARQL competency queries.
3. **Automation first** – use `go run ./cmd/bhashctl …` to orchestrate validation; legacy Python scripts are archival only.
4. **Semantic versioning** – ontology releases follow semver; changelogs must capture class/property additions and deprecations.
5. **Pull request checklist** – updated OWL/Turtle files, documentation, mapping tables, and validation evidence (SPARQL outputs or SHACL reports).

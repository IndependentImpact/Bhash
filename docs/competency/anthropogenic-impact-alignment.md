# Anthropogenic Impact Alignment Competency Answers

Phase 4 introduces the cross-ontology assets required to exercise CQ-IMPACT-001 and demonstrate AIAO-aligned reporting with Hedera data. The alignment sample now covers multiple SDG goals and mitigation schedules so stakeholder reviews can inspect diversified telemetry flows.

## CQ-IMPACT-001 – How do scheduled sustainability transactions translate into AIAO impact claims with supporting KPIs and telemetry provenance?

* **Status:** Validated via example graph and regression query in Phase 4 completion run.
* **Scope:** Show how scheduled mitigation actions funded on Hedera surface as AIAO impact claims with provenance to mirror datasets and sustainability KPIs.
* **Inputs:**
  * Ontology: `ontology/src/alignment/aiao.ttl`, `ontology/src/alignment/claimont.ttl`, `ontology/src/alignment/impactont.ttl`, `ontology/src/alignment/infocomm.ttl`
  * Sample graph: `ontology/examples/alignment-impact.ttl`
  * Query: `tests/queries/cq-impact-001.rq`

### Execution notes

1. Load the core and alignment ontologies alongside the anthropogenic impact sample graph into an RDF store.
2. Execute the SPARQL query to retrieve each `aiao:ImpactClaim`, its Hedera subject, associated KPI metrics, optional SDG goal, related scheduled mitigation action, and the communication node exposing evidence datasets.
3. Review the results to verify that Hedera concepts inherit semantics from AIAO, ClaimOnt, ImpactOnt, and InfoComm through the bridge axioms; confirm the query output matches `tests/fixtures/results/cq-impact-001.csv`.

### Sample result (derived from `ontology/examples/alignment-impact.ttl`)

| assertion | subject | indicator | value | sdgGoal | schedule | node |
| --------- | ------- | --------- | ----- | ------- | -------- | ---- |
| `ex:ReforestationAssertion` | `ex:CommunityTreasury` | `ex:ReforestationMetric` | `"1250"^^xsd:decimal` | `impactont:SDG15LifeOnLand` | `ex:ReforestationSchedule` | `ex:MirrorNodeAlpha` |
| `ex:SoilCarbonAssertion` | `ex:CommunityTreasury` | `ex:SoilCarbonMetric` | `"980"^^xsd:decimal` | `impactont:SDG13ClimateAction` | `ex:SoilCarbonSchedule` | `ex:MirrorNodeBeta` |

### Validation summary

* Regression query `tests/queries/cq-impact-001.rq` returns the expected assertion-to-dataset trace when executed against `ontology/examples/alignment-impact.ttl`.
* SHACL automation passes for the anthropogenic impact graph, confirming datatype constraints on `hedera:hasMetricValue`, `hedera:hasImpactVariance`, and required linkages to schedules and mirror datasets.
* Alignment blueprint and README now document the cross-ontology coverage so anthropogenic impact stakeholders can onboard without additional onboarding notes.

### Triple store data pilot

Run `python scripts/run_phase4_pilot.py` to initialise an Oxigraph-backed triple store, load the alignment ontologies and expanded sample data, execute `cq-impact-001`, and persist evidence under `build/pilots/phase4/`. A Go-based pilot command is planned as part of the Python helper retirement workstream.

* `build/pilots/phase4/cq-impact-001.csv` – live query results exported from the Oxigraph store.
* `build/pilots/phase4/shacl-report.txt` / `.ttl` – SHACL validation transcript derived from the triple store dump.
* `build/pilots/phase4/pilot-summary.json` – execution metadata (dataset counts, validation status, runtime in seconds).

### Next steps

* Extend sample data with adaptation-focused KPIs (e.g., flood resilience) to increase coverage across ClaimOnt mitigation categories.
* Wire live telemetry ingestion into the alignment pipeline so mirror datasets refresh automatically before stakeholder review workshops.

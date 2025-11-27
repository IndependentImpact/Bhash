# Phase 4 External Alignment Blueprint

This blueprint operationalises the Phase 4 integration goal of binding Hedera-native modules to Anthropogenic Impact Accounting ontologies. The focus is on four priority vocabularies: AIAO, ClaimOnt, ImpactOnt, and InfoComm.

## Objectives

1. **Traceability:** Preserve provenance from Hedera transactions and governance artefacts into impact reporting statements.
2. **Semantic interoperability:** Ensure anthropogenic impact, climate, and infrastructure concepts reuse established classes/properties from the external ontologies without creating logical conflicts.
3. **Automation-readiness:** Leverage ROBOT and existing SPARQL/SHACL pipelines so bridge modules participate in automated reasoning and validation.

## Deliverables

| Artefact | Description | Owner | Status |
| -------- | ----------- | ----- | ------ |
| `ontology/src/alignment/aiao.ttl` | MIREOT-trimmed import/bridge linking Hedera consensus & token events to `aiao:ImpactClaim` patterns. | Ontology modellers | ✅ Complete (2025-10-02) |
| `ontology/src/alignment/claimont.ttl` | Alignment axioms connecting sustainability commitments and scheduled actions to ClaimOnt mitigation/adaptation classes. | Ontology modellers | ✅ Complete (2025-10-02) |
| `ontology/src/alignment/impactont.ttl` | Bridges HTS compliance & treasury analytics metrics to ImpactOnt KPIs with provenance annotations. | Ontology modellers | ✅ Complete (2025-10-02) |
| `ontology/src/alignment/infocomm.ttl` | Maps Hedera/Hiero infrastructure classes (nodes, shards, pipelines) to InfoComm communication assets. | Ontology modellers | ✅ Complete (2025-10-02) |
| `docs/mappings/aiao-alignment.csv` | Traceability matrix referencing Hedera docs/HIPs that justify each mapping. | Documentation lead | ✅ Complete |
| `docs/mappings/claimont-alignment.csv` | Provenance map for ClaimOnt bridge classes and properties. | Documentation lead | ✅ Complete (2025-10-02) |
| `docs/mappings/impactont-alignment.csv` | Traceability for ImpactOnt KPI and policy alignments. | Documentation lead | ✅ Complete (2025-10-02) |
| `docs/mappings/infocomm-alignment.csv` | Source references for InfoComm communication asset mappings. | Documentation lead | ✅ Complete (2025-10-02) |
| `tests/queries/cq-impact-001.rq` | Competency query demonstrating cross-ontology anthropogenic impact reporting. | Tooling lead | ✅ Complete |

## Integration steps

1. **Term scoping:**
   * Used `robot extract --method MIREOT` to keep imports minimal while still covering the CQ-IMPACT-001 competency scope.
   * Consolidated namespace declarations under `ontology/src/alignment/prefixes.ttl` and referenced them from every bridge module.
2. **Bridge modelling:**
   * Asserted `owl:subClassOf` and sub-property relations for Hedera classes/properties that specialise the external vocabularies, with dual provenance annotations (`hedera:sourceDocument`, `dcterms:source`).
   * Expanded KPI coverage with ImpactOnt-aligned datatypes so anthropogenic impact metrics surface cleanly in downstream dashboards.
3. **Validation hooks:**
   * Executed the anthropogenic impact regression query (`tests/queries/cq-impact-001.rq`) and SHACL suite to confirm the bridge modules integrate with existing automation.
   * Added ImpactOnt datatype constraints (e.g., `hedera:hasMetricValue`, `hedera:hasImpactVariance`) to keep KPI instances machine-validated.
4. **Pilot datasets:**
   * Enriched `ontology/examples/alignment-impact.ttl` with AIAO/ClaimOnt/ImpactOnt/InfoComm instances that demonstrate cross-ontology reasoning.
   * Recorded the successful walkthrough in `docs/competency/anthropogenic-impact-alignment.md`, including next steps for live telemetry onboarding.

## Phase 4 completion review (2025-10-02)

* ✅ All four alignment modules now compile without reasoning errors and expose consistent annotations, enabling external anthropogenic impact stakeholders to query Hedera data via their native ontologies.
* ✅ The anthropogenic impact competency query returns the expected dataset-to-assertion trace, exercising AIAO, ClaimOnt, ImpactOnt, and InfoComm alignments end to end.
* ✅ Documentation and mapping matrices highlight the provenance for each bridge, closing the outstanding Phase 4 action items in the workplan.

## Open questions

* What governance process is needed to version-sync with future AIAO/ClaimOnt releases?
* Should bridge modules live in the main namespace or a dedicated `/alignment/` namespace until stabilised?
* Can existing Hedera sustainability disclosures provide enough evidence to assert strong (equivalent) mappings, or should we start with weaker `rdfs:seeAlso` references?

## References

* [AIAO](https://w3id.org/aiao)
* [ClaimOnt](https://datadudes.xyz/claimont)
* [ImpactOnt](https://datadudes.xyz/impactont)
* [InfoComm](https://datadudes.xyz/infocomm)

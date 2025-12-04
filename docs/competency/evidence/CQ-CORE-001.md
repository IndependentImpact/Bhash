# CQ-CORE-001 Evidence Bundle

**Question.** Which entities participate in consensus on mainnet?

## Documentation anchors

* Hedera mainnet consensus nodes are permissioned and run by the Governing Council, establishing the governance scope for
  onboarding. [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md)
* HIP-840 describes the Hiero roadmap and validator onboarding expectations as the network moves toward community validators.
  [HIP-840 Hiero network](https://hips.hedera.com/hip/hip-840)
* The transaction lifecycle and record exports detail how consensus participation materialises in observable artefacts and mirror
  datasets. [Transaction lifecycle overview](https://docs.hedera.com/hedera/core-concepts/transactions/transaction-lifecycle)

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| `ontology/examples/core-consensus.ttl` | Provides a worked example with mainnet, council, validator, staking accounts, transaction phases, and mirror dataset. | Repository sample (this commit). |
| `tests/queries/cq-core-001.rq` | SPARQL query that enumerates consensus participants and their governance mandate. | Repository query (this commit). |
| Mirror node REST `GET /api/v1/network/nodes` | Production endpoint for validator metadata when extending beyond the sample graph. | [Mirror node REST API](https://docs.hedera.com/hedera/mirror-node/sdks-and-apis/rest-api#tag/Network) |

## Modelling updates

1. Core ontology gains classes for `hedera:Account`, `hedera:ValidatorOnboardingProcess`, `hedera:TransactionPhase`, and
   `hedera:MirrorDataset`, along with properties (`hedera:hasAccount`, `hedera:targetsValidator`, `hedera:hasTransactionPhase`,
   `hedera:hasMandate`) to connect governance, identity, and lifecycle semantics.
2. Example data instantiates the onboarding decision, staking accounts, and transaction phases so the competency question can be
   executed without external dependencies.
3. Glossary entries in `docs/inventory/accounts-governance-glossary.md` and `docs/inventory/transaction-lifecycle.md` capture the
   underlying terminology and documentation links.

## Validation approach

* **SPARQL query:** `tests/queries/cq-core-001.rq`
  ```sparql
  PREFIX hedera: <https://hashgraphontology.xyz/core/>
  SELECT ?network ?validator ?operatorAccount ?mandate
  WHERE {
    ?network a hedera:Network ;
             hedera:isGovernedBy ?council .
    ?council hedera:hasMandate ?mandate .
    ?onboarding a hedera:ValidatorOnboardingProcess ;
               hedera:occursOn ?network ;
               hedera:targetsValidator ?validator .
    ?validator a hedera:ConsensusNode ;
               hedera:hasAccount ?operatorAccount ;
               hedera:participatesIn ?onboarding .
    ?operatorAccount hedera:isAccountOf ?validator .
  }
  ORDER BY ?network ?validator
  ```
* **Expected result (using the sample graph):** one row linking `ex:Mainnet`, `ex:Node3`, `ex:Node3Account`, and the HIP-840
  mandate string.
* **Future extension:** replace the sample dataset with live mirror-node exports to validate real validator rosters and ensure
  onboarding mandates reference published council records.

## Outstanding tasks

* Capture additional validators in the example dataset once HIP-840 guidance on community validators is publicly available.
* Add SHACL shapes ensuring every `hedera:ConsensusNode` has at least one `hedera:hasAccount` and associated `hedera:stakesTo`
  link before running against production data.
* Integrate the query into the automation backlog (AUT-005) so it runs automatically when ontology changes are proposed.

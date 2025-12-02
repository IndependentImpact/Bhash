# CQ-COMP-003 Evidence Bundle

**Question.** Which tokens classified as stablecoins (HIP-540) enforce KYC and freeze keys, and who controls those keys?

## Documentation anchors

* HIP-540 enumerates administrative keys (admin, KYC, freeze, wipe, metadata) and the mandate to remove or rotate them to guarantee immutability for compliant tokens. [HIP-540 specification](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md)
* HTS token creation guidance lists the security-sensitive keys (admin, supply, KYC, freeze, pause) and their operational responsibilities. [HTS token creation guide](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/tokens/hedera-token-service-hts-native-tokenization/token-creation.md)
* Mirror node REST `GET /api/v1/tokens/{id}` surfaces the freeze and KYC key metadata required to evaluate enforcement. [Token REST reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/token.md)

## Data sources

| Asset | Purpose | Access method |
| ----- | ------- | ------------- |
| HIP-540 dataset | Provides canonical list of stablecoin compliance requirements and expected key behaviours. | [HIP-540 specification](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md) |
| Mirror node `GET /api/v1/tokens` | Enumerates candidate tokens filtered by treasury or metadata indicating HIP-540 compliance. | [Token REST reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/token.md) |
| Mirror node `GET /api/v1/tokens/{id}` | Returns per-token key assignments (admin, freeze, kyc, fee schedule) for controller discovery. | [Token REST reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/token.md) |
| Mirror node `GET /api/v1/accounts/{id}` | Identifies controlling accounts for keys (e.g., treasury, admin key holder) and links to governance roles. | [Accounts REST reference](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/sdks-and-apis/rest-api/accounts.md) |

## Proposed modelling updates

1. Introduce `hedera:StablecoinToken` subclass tagged with `dcterms:identifier "HIP-540"` and connect to controller actors via `hedera:hasKeyAssignment` artefacts.
2. Represent key assets as `hedera:TokenKey` subclasses (`hedera:KYCKey`, `hedera:FreezeKey`) linked through `hedera:TokenKeyAssignment` nodes that record controller actors and governance roles.
3. Capture provenance between HIP-540 requirements and token instances using `hedera:provTrace` so compliance checks can assert traceability.

## Validation approach

* **SPARQL query prototype:**
  ```sparql
  PREFIX hedera: <https://hashgraphontology.xyx/hedera/core/>
  SELECT ?token ?treasury ?kycController ?freezeController
  WHERE {
    ?token a hedera:StablecoinToken ;
           hedera:hasTreasury ?treasury ;
           hedera:hasKeyAssignment ?kycAssignment , ?freezeAssignment .

    ?kycAssignment a hedera:KYCKeyAssignment ;
                   hedera:isControlledBy ?kycController .
    ?kycController hedera:hasRole ?kycRole .
    ?kycRole a hedera:KYCControllerRole .

    ?freezeAssignment a hedera:FreezeKeyAssignment ;
                      hedera:isControlledBy ?freezeController .
    ?freezeController hedera:hasRole ?freezeRole .
    ?freezeRole a hedera:FreezeControllerRole .
  }
  ```
  This query lists HIP-540 stablecoins with explicit controller accounts for KYC and freeze keys.
* **SHACL shape:** Validate that every `hedera:StablecoinToken` has `hasKeyAssignment` nodes for both KYC and freeze keys, each linking to an `hedera:Actor` with a governance role (e.g., compliance steward).

## Outstanding tasks

* Define heuristics (metadata tags, treasury lists) to identify HIP-540-compliant stablecoins in mirror node data; coordinate with product teams for authoritative registry if available.
* Decide how to represent key removal versus unusable keys to reflect HIP-540's allowance for zero-key immutability.
* Extend ingestion scripts to pull controller account aliases and map them to organisation-level actors for reporting.

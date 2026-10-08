# Accounts, Governance, and Node Role Glossary

This glossary expands the Phase 1 content audit with dedicated coverage of Hedera account constructs, governing bodies, and network roles. Each entry captures modelling cues and authoritative documentation links so the ontology can represent permissions, decision flows, and validator onboarding states with evidence.

## Accounts and credentials

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Hedera account | On-ledger identity created via transaction submission and funded by an existing account; stores keys, balance, and configuration state. | [Account creation](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/account-creation.md) | Model as subclass of `hedera:Artefact` (D-0007), linked to its controlling actor via `hedera:isAccountOf`, with required relationships to key material, payer account, and auto-renew configuration. |
| Account key list | Composite key (threshold or list) that controls authorisation on an account and downstream artefacts. | [Account properties – Keys](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/account-properties.md#keys) | Represent as reusable credential structure linked to accounts, files, topics, and tokens for governance pattern reuse. |
| Auto-create alias | Derived public-key or EVM address alias that allows creation of accounts upon first transfer. | [Auto account creation](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/auto-account-creation.md) | Capture alias identifiers and link to provenance of initial transfer event. |
| Treasury account | Account designated as token treasury with authority over supply and fee distributions. | [Tokenisation on Hedera](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/tokens/tokenization-on-hedera.md) | Relate to token artefacts via `hedera:managesArtefact`; include constraints for custom fee collectors. |
| Staking metadata account | Account that configures staking preferences, reward decline flag, and delegation target. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-nodes) | Model stake preferences as properties enabling inference of validator support and eligibility for rewards. |
| Hollow account | Auto-created account that has an alias but no key material until it is explicitly updated by the owner. | [Account properties – Hollow account](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/account-properties.md#what-is-a-hollow-account) | Represent as an intermediate account state requiring a follow-on key-assignment event to become fully controllable. |
| Account alias (public key/EVM) | Secondary identifier derived from a public key or EVM address that can reference an account in transfers and queries. | [Account properties – Account alias](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/accounts/account-properties.md#account-alias) | Capture alias type to support reasoning over wallet interoperability and DID mappings. |
| Decline rewards flag | Boolean account setting indicating that an account participates in staking without receiving reward distributions. | [Stake HBAR – Staking reward distribution](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/stake-hbar.md#staking-reward-distribution) | Model as account-level governance control tied to CQ-COMP-003 compliance scenarios and validator reward eligibility. |

## Governance bodies and decision actors

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Hedera Governing Council | Collective that operates permissioned consensus nodes, sets staking policy, and approves network-wide changes. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md) · [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch) | Represent as `hedera:GovernanceBody` with decision events tied to staking reward votes and validator onboarding approvals. |
| CoinCom (Council finance committee) | Council committee responsible for proposing and voting on staking reward rate updates. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch) | Model as specialised governance role connected to reward-rate modification processes and HIP mandates. |
| Compliance steward | Role accountable for ensuring token administrators manage keys (KYC, freeze, wipe) in line with HIP-540 immutability guarantees. | [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md) | Instantiate as `hedera:Role` associated with token admin accounts and referenced in competency questions on compliance controls. |
| Validator onboarding steward | Council-designated actor responsible for reviewing validator candidates and ratifying onboarding decisions through formal votes. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch) | Link to `hedera:ValidatorOnboardingProcess` activities with provenance pointers to HIP-840 and council vote records. |

## Governance processes and decisions

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Validator onboarding process | Formal council-governed workflow that evaluates validator candidates, records voting outcomes, and transitions approved nodes into active duty. | [HIP-840](https://hips.hedera.com/hip/hip-840); [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#phase-iii-staking-rewards-program-launch) | Model as `hedera:ValidatorOnboardingProcess` linking governance bodies, candidate validators, and resulting `hedera:GovernanceDecision` events. |
| Governance decision | Canonical record of council or committee approval that ratifies a policy, onboarding, or fee schedule change. | [Hedera Council](https://hedera.com/council); [Council governance docs](https://docs.hedera.com/hedera/core-concepts/governance) | Represent as `hedera:GovernanceDecision` events generated by governance processes with HIP mandates encoded via `hedera:hasMandate`. |
| Council vote record | Evidence of council member participation and quorum for a decision, often summarised in public meeting minutes or HIP releases. | [Council governance docs](https://docs.hedera.com/hedera/core-concepts/governance); [HIP repository](https://github.com/hiero-ledger/hiero-improvement-proposals) | Capture as subclasses of `hedera:VoteEvent` with provenance to source documents so SPARQL queries can enumerate supporting votes. |

## Token administration roles

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Token admin key custodian | Entity controlling the admin key that can modify or remove other token keys via `TokenUpdate` operations. | [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md#abstract) | Model as an accountable role with signing authority and lifecycle events when privileges are relinquished. |
| Token freeze key custodian | Delegate empowered to freeze or unfreeze token accounts, affecting transfer permissions. | [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md#abstract) | Represent as a compliance-critical role tied to CQ-COMP-003 with SHACL rules for freeze state transitions. |
| Token metadata key steward | Role responsible for updating token metadata and ensuring immutability once keys are removed. | [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md#abstract) | Capture as a governance participant whose actions must be logged prior to metadata key retirement. |

## Node and validator roles

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Consensus node operator | Council-operated node that submits transactions, participates in consensus, and enforces throttles. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md) | Link operators to node instances, staking weight, and governance approvals for onboarding or retirement. |
| Mirror node operator | Permissionless participant replicating consensus history for query APIs. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md) | Model as data publisher providing DCAT datasets and provenance for downstream analytics. |
| Validator onboarding candidate | Prospective Hiero validator subject to council review and HIP-540-style governance controls. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-nodes) · [HIP-540](https://raw.githubusercontent.com/hiero-ledger/hiero-improvement-proposals/main/HIP/hip-540.md) | Capture lifecycle states (applied, approved, active) and link to evidence bundles for CQ-GOV-001. |
| Reward account | System account 0.0.800 that disburses staking rewards once council-defined thresholds are met. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-reward-account) | Represent as managed artefact with policy constraints on balances and reward distribution triggers. |
| Consensus node account | Account identifier (0.0.x) associated with a council-operated node entry in the mainnet address book. | [Mainnet consensus nodes](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/networks/mainnet/mainnet-nodes/README.md#mainnet-node-address-book) | Model as `hedera:NodeAccount` linked to infrastructure metadata and onboarding approvals. |
| Proxy staking participant | Account that stakes to another account which in turn delegates to a node, enabling indirect reward flows. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md#staking-rewards) | Represent indirect staking chains to support reasoning over shared reward distributions and decline flags. |

## Staking lifecycle

| Term | Definition | Source | Modelling considerations |
| ---- | ---------- | ------ | ------------------------ |
| Staking process | Sequence through which an account delegates stake, accrues reward eligibility, and updates preferences such as decline rewards. | [Staking program](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/staking.md) | Model as `hedera:StakingProcess` with participants (accounts, validators) and emitted `hedera:StakingEvent` records capturing reward epochs. |
| Staking event | Observable change in staking state such as delegation updates, reward distribution, or penalty enforcement. | [Stake HBAR distribution](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/stake-hbar.md#staking-reward-distribution) | Represent as `hedera:StakingEvent` generated by staking processes to support time-series analytics and compliance checks. |
| Decline rewards flag | Account-level toggle signalling participation without reward payout. | [Stake HBAR distribution](https://raw.githubusercontent.com/hashgraph/hedera-docs/main/core-concepts/staking/stake-hbar.md#staking-reward-distribution) | Capture via the `hedera:declinesStakingRewards` datatype property so competency queries can reason over voluntary opt-outs. |

These entries now provide the dedicated account and governance vocabulary required by Immediate Action 2. Subsequent ontology iterations can attach SKOS notes and provenance citations using `hedera:sourceDocument` annotations aligned with this glossary.

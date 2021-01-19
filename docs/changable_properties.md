## Properties
This document is about properties we keep for Hubble and to show if they are updated by smart contract events. 

| Object  | Property | Description |  Event to update? |  Changeable on smart contract? |
| ------------- | ------------- | ------------- | ------------- | ------------- |
| Validator  | Validator ID  | the index of validator in SKALE deployed smart contract | No  | No  |
| Validator  | Name  | validator name | N/A  | Yes  |
| Validator  | Description  | validator description | N/A  | Yes  |
| Validator  | Validator Address  | validator address on SKALE (Address represents the 20 byte address of an Ethereum account) | ValidatorAddressChanged  | Yes  |
| Validator  | Requested Address  | requested address on SKALE (Address represents the 20 byte address of an Ethereum account) | N/A for adding, ValidatorAddressChanged for deleting  | Yes  |
| Validator  | Fee Rate  | fee rate | No  | No  |
| Validator  | Registration Time  | registration time to network | No  | No  |
| Validator  | Minimum Delegation Amount  | minimum delegation amount i.e. MDA | N/A  | Yes  |
| Validator  | Accept New Requests  | shows whether validator accepts new requests or not | N/A  | Yes  |
| Validator  | Authorized  | shows whether validator is authorized or not | ValidatorWasEnabled, ValidatorWasDisabled  | Yes  |
| Validator  | Active Nodes  | number of active nodes attached to the validator | N/A  | Deleting node which affect this result is available  |
| Validator  | Linked Nodes  | number of all nodes attached to the validator | N/A  | Deleting node which affect this result is available  |
| Validator  | Staked  | total stake amount | DelegationProposed, DelegationAccepted, DelegationRequestCanceledByUser, UndelegationRequested   | Yes  |
| Node  | Node ID | the index of node in SKALE deployed smart contract | No  | No  |
| Node  | Name  | node name | No  | No  |
| Node  | IP  | node ip | No  | No  |
| Node  | Public IP  | node public ip | No  | No  |
| Node  | Port  | node port | No  | No  |
| Node  | Start Block  | starting block height on ETH mainnet | No  | No  |
| Node  | Next Reward Date  | next reward time | BountyReceived **, BountyGot **  | Yes  |
| Node  | Last Reward Date  | last reward time | BountyReceived **, BountyGot **  | Yes  |
| Node  | Finish Time  | finish time | N/A  | Yes  |
| Node  | Status  | node status | N/A  | Yes  |
| Node  | Validator ID  | validator Id on SKALE network | N/A  | Yes (not mounting to another validator but deleting)  |
| Delegation  | Delegation ID  | the index of delegation in SKALE deployed smart contract | No  | No  |
| Delegation  | Holder  | Address of the token holder (Address represents the 20 byte address of an Ethereum account.) | No  | No  |
| Delegation  | Validator ID  | the index of validator in SKALE deployed smart contract | No  | No  |
| Delegation  | Block Height  | Block number at ETH mainnet | No  | No  |
| Delegation  | Transaction Hash  | transaction where delegation updated ( Hash represents the 32 byte Keccak256 hash of arbitrary data) | No  | No  |
| Delegation  | Amount  | delegation amount SKL unit | No  | No  |
| Delegation  | Period  | The duration delegation as chosen by the delegator | No  | No  |
| Delegation  | Created  | Creation time at ETH mainnet | No  | No  |
| Delegation  | Started  | started  epoch | No  | No  |
| Delegation  | Finished  | finished  epoch | No  | No  |
| Delegation  | Info  | delegation information | No  | No  |
| Delegation  | Status  | delegation status | N/A  | N/A  |

** available but not implemented on this repo yet!
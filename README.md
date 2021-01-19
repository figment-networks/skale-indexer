# SKALE INDEXER

This repository contains SKALE indexer code that is responsible for tracking SKALE network ethereum calls.

> :warning: **This repository is still undergoing heavy active development**: We will not consider this stable before the next announcement, and any code may be considerably changed.

## Requrements
    - Configured and installed go (to compile)
    - Ethereum node access. Should be the archival node

## Actions

### Compilation
To compile this project from sources you need to have go 1.15.6+ installed.

```bash
    make build
```

### Generate
To generate both swagger documentation and test mocks simply run:
```bash
    make generate
```

### Running
To run this project you just need to run the compiled binary supplying either config.json file or environment variables.

```bash
    ./indexer
```

Default parameters are available in `.env.default` file in this repository.
Keep in mind that you need to have SKALE abi files for the indexer to run. The easiest way to get it is clone `github.com/skalenetwork/skale-network` repository, and point indexer to `skale-network/releases/mainnet/skale-manager`. It will digest all abi files from entire directory.

This service is prepared to work with archival node of ethereum. It needs it to fetch previous states of smart contract. It's also possible to make it working with regular ethereum node. For that you must set env variable `ETHEREUM_NODE_TYPE` to "recent" to get only the latest states. At that mode the state may not be consistent, but you'll be able to get information about latest state.

## Calls

You can find detailed description of endpoints in swagger file.
Current version (0.0.1) allows only to get a range of records by calling:

```
    GET localhost:8885/getLogs?from=10940000&to=10950000
```

Where `to` and `from` is a range of the ethereum blocks to get this information from.
This operation is indempotent, and should only update records in case of previous failures

## Properties


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
| Validator  | Pending  | ? | ?  | ?  |
| Validator  | Rewards  | ? | ?  | ?  |
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
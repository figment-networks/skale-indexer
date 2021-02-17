# SKALE INDEXER

This repository contains SKALE indexer code that is responsible for tracking SKALE network Ethereum calls.

## Requirements
    - Configured and installed go (to compile)
    - Ethereum node access (archive preferred for full functionality)

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
To run this project you just need to run the compiled binary supplying either the config.json file or environment variables.

```bash
    ./indexer
```

Default parameters are available in `.env.default` file in this repository.

Keep in mind that you need to have SKALE abi files for the indexer to run (`ABI_DIR`) . The easiest way to get it is clone `github.com/skalenetwork/skale-network` repository, and point indexer to `skale-network/releases/mainnet/skale-manager`. It will digest all abi files from entire directory.

This service is designed to work with Ethereum archive nodes, as it needs it to fetch the previous states of the smart contracts. 

It's also possible to configure the indexer to work with a regular Ethereum node. For that you must set env variable `ETHEREUM_NODE_TYPE` to "recent" to get only the latest states. With this mode the state may not be consistent, but you'll be able to get information about latest state only.

It if possible to run this service in read only mode. to enable that you just need to set `ENABLE_SCRAPER=true`

Skale indexer is using special kind of ABI deployment using proxy. This is why additional events that are not defined in abi files may arrive. Those additional events should be passed in file  `ADDITIONAL_ABI`. The file is available in this repository as `.additional.abi.json`.

Because ethereum blocks are different we need to declare "zero" block and time, after what we gonna start probing for events `ETHEREUM_SMALLEST_BLOCK_NUMBER` `ETHEREUM_SMALLEST_BLOCK_TIME`

## Calls

You can find detailed description of endpoints in swagger file.
To get a range of requests you need to call:

```
    GET localhost:8885/getLogs?from=10940000&to=10950000
```

Where `to` and `from` is a range of the ethereum blocks to get this information from.
This operation is indempotent, and should only update records in case of previous failures

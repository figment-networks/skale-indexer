# Change Log

## [0.0.9] - 2021-06-28
### Added
- Adds `limit` and `offset` params to `/system_events`
### Changed
- Modifies the `ORDER BY` lookup for the System Events to return latest block first instead of the oldest

## [0.0.8] - 2021-06-15
### Fixed
- bugfix invalid lookup for `/system_events?address=` due to no bigInt conversion

## [0.0.7] - 2021-05-18
### Added
### Changed
### Fixed

## [0.0.6] - 2021-05-18
### Fixed
- bugfix ambiguous column for in `/delegations` with timeline

## [0.0.5] - 2021-05-18
### Fixed
- bugfix ambiguous column for in `/delegations`

## [0.0.4] - 2021-04-10
### Added
- Return `validator_name` in `/delegations`
- Add query param `address` to `/validators`
- New index `idx_val_addr` on validators table

### Changed
- Various changes for hubble consistency
### Fixed

## [0.0.2] - 2021-02-17

First Pre-Hubble version

### Added
- Cache for Addresses and Delegations is now in place
- Integration with scheduler is added to the node to get new blocks periodically
- In case of abi change, we check for earlier versions of event in all previously loaded versions of abi
- There is now additional abi you can set on top of the indexer, it was needed for proxy certificate events
- Added code that fetches transactions after crossing month.

### Changed
- All data currently is taken from the node. We're no longer calculate total stake or linked/active nodes in the database.
- Added timestamps to some events for hubble
- Some additional parameters in http api

### Fixed
- Many minor and more serious bugs
- Postgres `deadlock detected` error on saving nodes


## [0.0.1] - 2021-01-14

Initial release

### Added
### Changed
### Fixed


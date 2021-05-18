# Change Log

## [0.0.5] - 2021-05-18
### Added
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


package webapi

import "time"

// EventParams a set of fields to be used for events (contract events) search
type EventParams struct {
	// Id - represents the bound id of the event
	//
	// supposed to used together with Type, i.e. required when Type is used
	// format: unsigned integer
	Id uint64
	// Type - filtering events by event type
	//
	// supposed to used together with Id, i.e. required when Id is used
	// example: "validator", "delegation", "node", "token"
	Type string
	// TimeFrom - the inclusive beginning of the time range for event time
	//
	// supposed to be sent with time to
	// required: true
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeFrom time.Time
	// TimeTo - the inclusive ending of the time range for event time
	//
	// supposed to be sent with time from
	// required: true
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeTo time.Time
}

// NodeParams a set of fields to be used for nodes search
type NodeParams struct {
	// NodeID - the index of node in SKALE deployed smart contract
	//
	// format: unsigned integer
	NodeID string `json:"id"`
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"validator_id"`
}

// AccountParams a set of fields to be used for accounts search
type AccountParams struct {
	// Type - type of the account
	//
	// example: "validator", "delegator", "default"
	Type string `json:"type"`
	// Address - account address i.e. holder
	//
	// format: hexadecimal
	Address string `json:"address"`
}

// DelegationParams a set of fields to be used for accounts search
type DelegationParams struct {
	// DelegationID - the index of delegation in SKALE deployed smart contract
	//
	// format: unsigned integer
	DelegationID string `json:"id"`
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"validator_id"`
	// TimeFrom - the inclusive beginning of the time range for delegation created time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeFrom time.Time `json:"from"`
	// TimeTo - the inclusive ending of the time range for delegation created time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeTo time.Time `json:"to"`
	// Timeline - returns whether the latest or delegation changes timeline
	//
	// case false to fetch recent info for filtered delegations
	// case true to fetch whole delegations for the filter
	Timeline bool `json:"timeline"`
}

// ValidatorParams a set of fields to be used for validators search
type ValidatorParams struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"id"`
	// TimeFrom - the inclusive beginning of the time range for registration time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeFrom time.Time `json:"from"`
	// TimeTo - the inclusive ending of the time range for registration time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeTo time.Time `json:"to"`
}

// ValidatorStatisticsParams a set of fields to be used for validator statistics search
type ValidatorStatisticsParams struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"id"`
	// Type - statistics type
	//
	// example: "TOTAL_STAKE", "ACTIVE_NODES" etc...
	Type string `json:"type"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64    `json:"height"`
	TimeFrom    time.Time `json:"from"`
	TimeTo      time.Time `json:"to"`
	// Timeline - returns whether the latest or statistics changes timeline
	//
	// case false to fetch recent info for filtered statistics
	// case true to fetch whole statistics for the filter
	Timeline bool `json:"timeline"`
}

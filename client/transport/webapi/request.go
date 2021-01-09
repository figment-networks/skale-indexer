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
	// TimeFrom - filtering events by event time from
	//
	// supposed to be sent with time to
	// required: true
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeFrom time.Time
	// TimeTo - filtering events by events time to
	//
	// supposed to be sent with time from
	// required: true
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeTo time.Time
}

// NodeParams a set of fields to be used for nodes search
type NodeParams struct {
	// NodeID - node Id on SKALE network
	//
	// format: unsigned integer
	NodeID string `json:"id"`
	// ValidatorID - node Id on SKALE network
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
	// DelegationID - delegation Id on SKALE network
	//
	// format: unsigned integer
	DelegationID string `json:"id"`
	// ValidatorID - delegation Id on SKALE network
	//
	// format: unsigned integer
	ValidatorID string `json:"validator_id"`
	// TimeFrom - filtering delegations by created time from
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeFrom time.Time `json:"from"`
	// TimeTo - filtering delegations by created time to
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeTo time.Time `json:"to"`
	// Timeline - fetching whether the latest or time chart for filtered delegations
	//
	// case false to fetch recent info for filtered delegations
	// case true to fetch whole delegations for the filter
	Timeline bool `json:"timeline"`
}

// ValidatorParams a set of fields to be used for validators search
type ValidatorParams struct {
	// ValidatorID - delegation Id on SKALE network
	//
	// format: unsigned integer
	ValidatorID string `json:"id"`
	// TimeFrom - filtering validators by registration time from
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeFrom time.Time `json:"from"`
	// TimeTo - filtering validators by registration time to
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: "2006-01-02T15:04:05Z07:00"
	TimeTo time.Time `json:"to"`
}

// ValidatorStatisticsParams a set of fields to be used for validator statistics search
type ValidatorStatisticsParams struct {
	// ValidatorID - delegation Id on SKALE network
	//
	// format: unsigned integer
	ValidatorID string `json:"id"`
	// Type - statistics type
	//
	// example: "TOTAL_STAKE", "ACTIVE_NODES"
	Type string `json:"type"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64    `json:"height"`
	TimeFrom    time.Time `json:"from"`
	TimeTo      time.Time `json:"to"`
	// Timeline - fetching whether the latest or time chart for filtered statistics
	//
	// case false to fetch recent info for filtered statistics
	// case true to fetch whole statistics for the filter
	Timeline bool `json:"timeline"`
}

package webapi

import "time"

// EventParams a set of fields to be used for events (contract events) search
// swagger:model
type EventParams struct {
	// Id - represents the bound id of the event
	//
	// supposed to used together with Type, i.e. required when Type is used
	// format: unsigned integer
	// required: false
	ID uint64 `json:"id,omitempty"`
	// Type - filtering events by event type
	//
	// supposed to used together with Id, i.e. required when Id is used
	// example: validator
	Type string `json:"type"`
	// TimeFrom - the inclusive beginning of the time range for event time
	//
	// supposed to be sent with time to
	// required: true
	// time format: RFC3339
	// example: 2020-09-22T12:42:31Z
	TimeFrom time.Time `json:"from"`
	// TimeTo - the inclusive ending of the time range for event time
	//
	// supposed to be sent with time from
	// required: true
	// time format: RFC3339
	// example: 2021-09-22T12:42:31Z
	TimeTo time.Time `json:"to"`
	// Limit - Limit of the records per page
	//
	// required: false
	Limit uint64 `json:"limit"`
	// Offset - Offset of the records in pages
	//
	// required: false
	Offset uint64 `json:"offset"`
}

// NodeParams a set of fields to be used for nodes search
// swagger:model
type NodeParams struct {
	// NodeID - the index of node in SKALE deployed smart contract
	//
	// format: unsigned integer
	NodeID string `json:"id"`
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"validator_id"`
	// Status - node status
	//
	// example: Active
	Status string `json:"status"`
	// Limit - Limit of the records per page
	//
	// required: false
	Limit uint64 `json:"limit"`
	// Offset - Offset of the records in pages
	//
	// required: false
	Offset uint64 `json:"offset"`
}

// AccountParams a set of fields to be used for accounts search
// swagger:model
type AccountParams struct {
	// Type - type of the account
	//
	Type string `json:"type"`
	// Address - account address i.e. holder
	//
	// format: hexadecimal
	Address string `json:"address"`
	// Limit - Limit of the records per page
	//
	// required: false
	Limit uint64 `json:"limit"`
	// Offset - Offset of the records in pages
	//
	// required: false
	Offset uint64 `json:"offset"`
}

// DelegationParams a set of fields to be used for accounts search
// swagger:model
type DelegationParams struct {
	// DelegationID - the index of delegation in SKALE deployed smart contract
	//
	// format: unsigned integer
	DelegationID string `json:"id"`
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"validator_id"`
	// Holder - holder address
	//
	// format: hexadecimal
	Holder string `json:"holder"`
	// TimeAt - point time for that the validation statuses will be returned
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: 2021-09-22T12:42:31Z
	TimeAt time.Time `json:"at"`
	// Timeline - returns whether the latest or delegation changes timeline
	//
	// case false to fetch recent info for filtered delegations
	// case true to fetch whole delegations for the filter
	Timeline bool `json:"timeline"`
	// Limit - Limit of the records per page
	//
	// required: false
	Limit uint64 `json:"limit"`
	// Offset - Offset of the records in pages
	//
	// required: false
	Offset uint64 `json:"offset"`
	// State - states to filter by
	//
	// required: false
	State []string `json:"state"`
}

// ValidatorParams a set of fields to be used for validators search
// swagger:model
type ValidatorParams struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"id"`
	// TimeFrom - the inclusive beginning of the time range for registration time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: 2020-09-22T12:42:31Z
	TimeFrom time.Time `json:"from"`
	// TimeTo - the inclusive ending of the time range for registration time
	//
	// supposed to be sent with time from
	// time format: RFC3339
	// example: 2021-09-22T12:42:31Z
	TimeTo time.Time `json:"to"`
	// Authorized - get 1 = only authorized, 2 = only not authorized
	//
	// example: 1
	Authorized uint8 `json:"authorized"`
	// Address - validator address
	//
	// format: hexadecimal
	Address string `json:"address"`
	// Limit - Limit of the records per page
	//
	// required: false
	Limit uint64 `json:"limit"`
	// Offset - Offset of the records in pages
	//
	// required: false
	Offset uint64 `json:"offset"`
}

// ValidatorStatisticsParams a set of fields to be used for validator statistics search
// swagger:model
type ValidatorStatisticsParams struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// format: unsigned integer
	ValidatorID string `json:"id"`
	// Type - statistics type
	//
	// example: TOTAL_STAKE
	Type string `json:"type"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"height"`
	// TimeFrom - the inclusive beginning of the time range for block time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: 2020-09-22T12:42:31Z
	TimeFrom time.Time `json:"from"`
	// TimeTo - the inclusive ending of the time range for block time
	//
	// supposed to be sent with time to
	// time format: RFC3339
	// example: 2021-09-22T12:42:31Z
	TimeTo time.Time `json:"to"`
	// Timeline - returns whether the latest or statistics changes timeline
	//
	// case false to fetch recent info for filtered statistics
	// case true to fetch whole statistics for the filter
	Timeline bool `json:"timeline"`
	// Limit - Limit of the records per page
	//
	// required: false
	Limit uint64 `json:"limit"`
	// Offset - Offset of the records in pages
	//
	// required: false
	Offset uint64 `json:"offset"`
}

// SystemEventParams a set of fields to be used for system events
type SystemEventParams struct {
	After       uint64 `json:"after"`
	Kind        string `json:"kind"`
	Address     string `json:"address"`
	ValidatorID string `json:"validator_id"`
	ID          string `json:"id"`
	SenderID    uint64 `json:"sender_id"`
	ReceiverID  uint64 `json:"receiver_id"`
	Limit       uint64 `json:"limit"`
	Offset      uint64 `json:"offset"`
}

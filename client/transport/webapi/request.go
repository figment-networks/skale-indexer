package webapi

import "time"

type EventParams struct {
	RecordId string
	Id       uint64
	Type     string
	TimeFrom time.Time
	TimeTo   time.Time
}

type NodeParams struct {
	NodeID      string `json:"id"`
	ValidatorID string `json:"validator_id"`
}

type AccountParams struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type DelegationParams struct {
	DelegationID string    `json:"id"`
	ValidatorID  string    `json:"validator_id"`
	TimeFrom     time.Time `json:"from"`
	TimeTo       time.Time `json:"to"`
	Timeline     bool      `json:"timeline"`
}

type ValidatorParams struct {
	ValidatorID string    `json:"id"`
	TimeFrom    time.Time `json:"from"`
	TimeTo      time.Time `json:"to"`
}

type ValidatorStatisticsParams struct {
	ValidatorID string    `json:"id"`
	Type        string    `json:"type"`
	BlockHeight uint64    `json:"height"`
	TimeFrom    time.Time `json:"from"`
	TimeTo      time.Time `json:"to"`
	Timeline    bool      `json:"timeline"`
}

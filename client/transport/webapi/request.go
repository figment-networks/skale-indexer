package webapi

import "time"

type ValidatorParams struct {
	ValidatorID string    `json:"id"`
	TimeFrom    time.Time `json:"from"`
	TimeTo      time.Time `json:"to"`
}

type ValidatorStatisticsParams struct {
	ValidatorID      string    `json:"id"`
	StatisticsTypeVS string    `json:"type"`
	BlockHeight      uint64    `json:"height"`
	Timeline         bool      `json:"timeline"`
	TimeFrom         time.Time `json:"from"`
	TimeTo           time.Time `json:"to"`
}

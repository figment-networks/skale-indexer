package structs

import (
	"time"
)

const Layout = time.RFC3339

type EventParams struct {
	RecordId string
	Id       uint64
	Type     string
	TimeFrom time.Time
	TimeTo   time.Time
}

type DelegationParams struct {
	ValidatorId  string
	DelegationId string
	BlockHeight  string
	Holder       string
	TimeFrom     time.Time
	TimeTo       time.Time
}

type NodeParams struct {
	NodeId      string
	ValidatorId string
}

type AccountParams struct {
	Type    string
	Address string
}

type ValidatorParams struct {
	ValidatorId string
	TimeFrom    time.Time
	TimeTo      time.Time
}

type ValidatorStatisticsParams struct {
	ValidatorId      string
	StatisticsTypeVS string
	BlockHeight      uint64
	TimeFrom         time.Time
	TimeTo           time.Time
}

type DelegatorStatisticsParams struct {
	Holder           string
	StatisticsTypeDS string
	BlockHeight      uint64
	TimeFrom         time.Time
	TimeTo           time.Time
}

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
	ValidatorID string
	TimeFrom    time.Time
	TimeTo      time.Time
}

type ValidatorStatisticsParams struct {
	ValidatorID      string
	StatisticsTypeVS StatisticTypeVS
	BlockHeight      uint64
	Timeline         bool
	TimeFrom         time.Time
	TimeTo           time.Time
}

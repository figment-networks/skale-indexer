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
	ValidatorId    string
	DelegationId   string
	ETHBlockHeight string
	Recent         bool
	TimeFrom       time.Time
	TimeTo         time.Time
}

type NodeParams struct {
	ValidatorId string
	Recent      bool
}

type ValidatorParams struct {
	ValidatorId string
	Recent      bool
	TimeFrom    time.Time
	TimeTo      time.Time
}

type AccountParams struct {
	Type    string
	Address bool
}

type ValidatorStatisticsParams struct {
	ValidatorId     string
	StatisticTypeVS string
	BlockHeight     uint64
	Recent          bool
	TimeFrom        time.Time
	TimeTo          time.Time
}

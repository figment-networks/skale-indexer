package structs

import (
	"time"
)

const Layout = time.RFC3339

type EventParams struct {
	Id       uint64
	Type     string
	TimeFrom time.Time
	TimeTo   time.Time
}

type DelegationParams struct {
	ValidatorID  string
	DelegationID string
	TimeFrom     time.Time
	TimeTo       time.Time
}

type NodeParams struct {
	NodeID      string
	ValidatorID string
	Status      string
	Address     string
}

type AccountParams struct {
	Type    string
	Address string
}

type ValidatorParams struct {
	ValidatorID    string
	OrderBy        string
	OrderDirection string
	TimeFrom       time.Time
	TimeTo         time.Time
}

type ValidatorStatisticsParams struct {
	ValidatorID string
	Type        StatisticTypeVS
	BlockHeight uint64
	Time        time.Time
	TimeFrom    time.Time
	TimeTo      time.Time
}
type SystemEventParams struct {
	After      uint64
	Kind       string
	Address    string
	SenderID   uint64
	ReceiverID uint64
}

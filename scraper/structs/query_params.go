package structs

import (
	"time"
)

const Layout = time.RFC3339

type ThreeState uint8

const (
	StateNotSet ThreeState = iota
	StateTrue
	StateFalse
)

type EventParams struct {
	Id       uint64
	Type     string
	TimeFrom time.Time
	TimeTo   time.Time

	Limit  uint64
	Offset uint64
}

type DelegationParams struct {
	ValidatorID  string
	DelegationID string
	Holder       string
	State        []DelegationState
	TimeAt       time.Time

	TimeFrom time.Time
	TimeTo   time.Time

	Limit  uint64
	Offset uint64
}

type NodeParams struct {
	NodeID      string
	ValidatorID string
	Status      string
	Address     string

	Limit  uint64
	Offset uint64
}

type AccountParams struct {
	Type    string
	Address string

	Limit  uint64
	Offset uint64
}

type ValidatorParams struct {
	ValidatorID    string
	OrderBy        string
	OrderDirection string
	Authorized     ThreeState
	TimeFrom       time.Time
	TimeTo         time.Time
	Address        string

	Limit  uint64
	Offset uint64
}

type ValidatorStatisticsParams struct {
	ValidatorID string
	Type        StatisticTypeVS
	BlockHeight uint64
	BlockTime   time.Time
	TimeFrom    time.Time
	TimeTo      time.Time

	Limit  uint64
	Offset uint64
}

type SystemEventParams struct {
	After      uint64
	Kind       string
	Address    string
	SenderID   uint64
	ReceiverID uint64

	Limit 		 uint64
	Offset 		 uint64
}

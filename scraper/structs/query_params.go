package structs

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const Layout = time.RFC3339

type QueryParams struct {
	Id              string
	ValidatorId     uint64
	DelegationId    uint64
	Recent          bool
	Holder          common.Address
	ETHBlockHeight  uint64
	StatisticTypeVS StatisticTypeVS
	BoundType       string
	BoundId         []uint64
	TimeFrom        time.Time
	TimeTo          time.Time
}

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
	Kind    string
	Id      string
	Address bool
	Recent  bool
}

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

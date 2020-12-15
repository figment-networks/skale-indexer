package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
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
	BoundId     	[]uint64
	TimeFrom        time.Time
	TimeTo          time.Time
}

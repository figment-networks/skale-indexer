package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type QueryParams struct {
	Id              string
	ValidatorId     uint64
	Recent          bool
	Holder          common.Address
	ETHBlockHeight  uint64
	StatisticTypeVS StatisticTypeVS
	TimeFrom        time.Time
	TimeTo          time.Time
}

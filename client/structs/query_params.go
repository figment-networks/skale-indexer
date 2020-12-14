package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type QueryParams struct {
	Id              string
	ValidatorId     *big.Int
	Recent          bool
	Holder          common.Address
	ETHBlockHeight  uint64
	StatisticTypeVS StatisticTypeVS
	TimeFrom        time.Time
	TimeTo          time.Time
}

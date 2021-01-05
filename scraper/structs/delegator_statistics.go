package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type DelegatorStatistics struct {
	ID               string           `json:"id"`
	CreatedAt        time.Time        `json:"created_at"`
	Holder           common.Address   `json:"holder"`
	Amount           *big.Int         `json:"amount"`
	BlockHeight      uint64           `json:"block_height"`
	StatisticsTypeDS StatisticsTypeDS `json:"statistics_type"`
}

type StatisticsTypeDS int

const (
	DelegatorStatisticsTypeClaimedRewards StatisticsTypeDS = iota + 1
	DelegatorStatisticsTypeUnclaimedRewards
)

func (k StatisticsTypeDS) String() string {
	switch k {
	case DelegatorStatisticsTypeClaimedRewards:
		return "CLAIMED_REWARDS"
	case DelegatorStatisticsTypeUnclaimedRewards:
		return "UNCLAIMED_REWARDS"
	default:
		return "unknown"
	}
}

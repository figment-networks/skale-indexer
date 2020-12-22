package structs

import (
	"math/big"
	"time"
)

type ValidatorStatistics struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ValidatorId   *big.Int        `json:"validator_id"`
	Amount        uint64          `json:"amount"`
	BlockHeight   uint64          `json:"block_height"`
	StatisticType StatisticTypeVS `json:"statistics_type"`
}

type StatisticTypeVS int

const (
	ValidatorStatisticsTypeTotalStake StatisticTypeVS = iota + 1
	ValidatorStatisticsTypeActiveNodes
	ValidatorStatisticsTypeLinkedNodes
)

func (k StatisticTypeVS) String() string {
	switch k {
	case ValidatorStatisticsTypeTotalStake:
		return "TOTAL_STAKE"
	case ValidatorStatisticsTypeActiveNodes:
		return "ACTIVE_NODES"
	case ValidatorStatisticsTypeLinkedNodes:
		return "LINKED_NODES"
	default:
		return "unknown"
	}
}

package structs

import (
	"math/big"
	"time"
)

type ValidatorStatistics struct {
	ID             string           `json:"id"`
	CreatedAt      time.Time        `json:"created_at"`
	ValidatorId    *big.Int         `json:"validator_id"`
	Amount         *big.Int         `json:"amount"`
	BlockHeight    uint64           `json:"block_height"`
	StatisticsType StatisticsTypeVS `json:"statistics_type"`
}

type StatisticsTypeVS int

const (
	ValidatorStatisticsTypeTotalStake StatisticsTypeVS = iota + 1
	ValidatorStatisticsTypeActiveNodes
	ValidatorStatisticsTypeLinkedNodes
	ValidatorStatisticsTypeUnclaimedRewards
	ValidatorStatisticsTypeClaimedRewards
	ValidatorStatisticsTypeBounty
)

func (k StatisticsTypeVS) String() string {
	switch k {
	case ValidatorStatisticsTypeTotalStake:
		return "TOTAL_STAKE"
	case ValidatorStatisticsTypeActiveNodes:
		return "ACTIVE_NODES"
	case ValidatorStatisticsTypeLinkedNodes:
		return "LINKED_NODES"
	case ValidatorStatisticsTypeUnclaimedRewards:
		return "UNCLAIMED_REWARDS"
	case ValidatorStatisticsTypeClaimedRewards:
		return "CLAIMED_REWARDS"
	case ValidatorStatisticsTypeBounty:
		return "BOUNTY"
	default:
		return "unknown"
	}
}

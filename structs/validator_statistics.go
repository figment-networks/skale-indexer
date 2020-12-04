package structs

import (
	"time"
)

type ValidatorStatistics struct {
	ID            string           `json:"id"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	Status        DelegationStatus `json:"status"`
	ValidatorId   uint64           `json:"validator_id"`
	Amount        uint64           `json:"amount"`
	StatisticType StatisticTypeVS  `json:"statistics_type"`
}

type StatisticTypeVS int

const (
	TotalStakeStatisticsTypeVS StatisticTypeVS = iota + 1
	ActiveNodesStatisticsTypeVS
	LinkedNodesStatisticsTypeVS
)

func (k StatisticTypeVS) String() string {
	switch k {
	case TotalStakeStatisticsTypeVS:
		return "TOTAL_STAKE"
	case ActiveNodesStatisticsTypeVS:
		return "ACTIVE_NODES"
	case LinkedNodesStatisticsTypeVS:
		return "LINKED_NODES"
	default:
		return "unknown"
	}
}

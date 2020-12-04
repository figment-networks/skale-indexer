package structs

import (
	"time"
)

type DelegationStatistics struct {
	ID            string           `json:"id"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	Status        DelegationStatus `json:"status"`
	ValidatorId   uint64           `json:"validator_id"`
	Amount        uint64           `json:"amount"`
	StatisticType StatisticTypeDS  `json:"statistics_type"`
}

type StatisticTypeDS int

const (
	StatesStatisticsTypeDS StatisticTypeDS = iota + 1
	NextEpochStatisticsTypeDS
)

func (k StatisticTypeDS) String() string {
	switch k {
	case StatesStatisticsTypeDS:
		return "STATE_STATISTICS"
	case NextEpochStatisticsTypeDS:
		return "NEXT_EPOCH_STATISTICS"
	default:
		return "unknown"
	}
}

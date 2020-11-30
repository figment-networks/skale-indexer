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
	StatisticType StatisticType    `json:"statistics_type"`
}

type StatisticType int

const (
	StatesStatisticsType StatisticType = iota + 1
	NextEpochStatisticsType
)

func (k StatisticType) String() string {
	switch k {
	case StatesStatisticsType:
		return "STATE_STATISTICS"
	case NextEpochStatisticsType:
		return "NEXT_EPOCH_STATISTICS"
	default:
		return "unknown"
	}
}

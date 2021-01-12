package structs

import (
	"math/big"
	"time"
)

type ValidatorStatistics struct {
	ID          string          `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	ValidatorID *big.Int        `json:"validator_id"`
	Amount      *big.Int        `json:"amount"`
	BlockHeight uint64          `json:"block_height"`
	Type        StatisticTypeVS `json:"type"`
}

const (
	ValidatorStatisticsTypeTotalStake StatisticTypeVS = iota + 1
	ValidatorStatisticsTypeActiveNodes
	ValidatorStatisticsTypeLinkedNodes
	ValidatorStatisticsTypeMDR
	ValidatorStatisticsTypeFee
	ValidatorStatisticsTypeAuthorized
	ValidatorStatisticsTypeValidatorAddress
	ValidatorStatisticsTypeRequestedAddress
)

var (
	StatisticTypes = map[string]StatisticTypeVS{
		"TOTAL_STAKE":       ValidatorStatisticsTypeTotalStake,
		"ACTIVE_NODES":      ValidatorStatisticsTypeActiveNodes,
		"LINKED_NODES":      ValidatorStatisticsTypeLinkedNodes,
		"MDR":               ValidatorStatisticsTypeMDR,
		"FEE":               ValidatorStatisticsTypeFee,
		"AUTHORIZED":        ValidatorStatisticsTypeAuthorized,
		"VALIDATOR_ADDRESS": ValidatorStatisticsTypeValidatorAddress,
		"REQUESTED_ADDRESS": ValidatorStatisticsTypeRequestedAddress,
	}
)

func GetTypeForValidatorStatistics(s string) (StatisticTypeVS, bool) {
	t, ok := StatisticTypes[s]
	return t, ok
}

type StatisticTypeVS uint

func (k StatisticTypeVS) String() string {
	switch k {
	case ValidatorStatisticsTypeTotalStake:
		return "TOTAL_STAKE"
	case ValidatorStatisticsTypeActiveNodes:
		return "ACTIVE_NODES"
	case ValidatorStatisticsTypeLinkedNodes:
		return "LINKED_NODES"
	case ValidatorStatisticsTypeMDR:
		return "MDR"
	case ValidatorStatisticsTypeFee:
		return "FEE"
	case ValidatorStatisticsTypeAuthorized:
		return "AUTHORIZED"
	case ValidatorStatisticsTypeValidatorAddress:
		return "VALIDATOR_ADDRESS"
	case ValidatorStatisticsTypeRequestedAddress:
		return "REQUESTED_ADDRESS"
	default:
		return "unknown"
	}
}

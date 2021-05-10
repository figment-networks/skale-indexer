package structs

import (
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
)

type Delegation struct {
	ID               uuid.UUID       `json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	DelegationID     *big.Int        `json:"delegation_id"`
	Holder           common.Address  `json:"holder"`
	ValidatorID      *big.Int        `json:"validatorId"`
	BlockHeight      uint64          `json:"block_height"`
	TransactionHash  common.Hash     `json:"transaction_hash"`
	Amount           *big.Int        `json:"amount"`
	DelegationPeriod *big.Int        `json:"delegationPeriod"`
	Created          time.Time       `json:"created"`
	Started          *big.Int        `json:"started"`
	Finished         *big.Int        `json:"finished"`
	Info             string          `json:"info"`
	State            DelegationState `json:"state"`

	// JOIN
	ValidatorName string `json:"validator_name"`
}

type DelegationState uint

const (
	DelegationStatePROPOSED DelegationState = iota
	DelegationStateACCEPTED
	DelegationStateCANCELED
	DelegationStateREJECTED
	DelegationStateDELEGATED
	DelegationStateUNDELEGATION_REQUESTED
	DelegationStateCOMPLETED

	DelegationStateUNKNOWN DelegationState = 666
)

func (k DelegationState) String() string {
	switch k {
	case DelegationStatePROPOSED:
		return "PROPOSED"
	case DelegationStateACCEPTED:
		return "ACCEPTED"
	case DelegationStateCANCELED:
		return "CANCELED"
	case DelegationStateREJECTED:
		return "REJECTED"
	case DelegationStateDELEGATED:
		return "DELEGATED"
	case DelegationStateUNDELEGATION_REQUESTED:
		return "UNDELEGATION_REQUESTED"
	case DelegationStateCOMPLETED:
		return "COMPLETED"
	default:
		return "UNKNOWN"
	}
}

func DelegationStateFromString(s string) DelegationState {
	switch s {
	case "PROPOSED":
		return DelegationStatePROPOSED
	case "ACCEPTED":
		return DelegationStateACCEPTED
	case "CANCELED":
		return DelegationStateCANCELED
	case "REJECTED":
		return DelegationStateREJECTED
	case "DELEGATED":
		return DelegationStateDELEGATED
	case "UNDELEGATION_REQUESTED":
		return DelegationStateUNDELEGATION_REQUESTED
	case "COMPLETED":
		return DelegationStateCOMPLETED
	default:
		return DelegationStateUNKNOWN
	}
}

type DelegationSummary struct {
	Count  *big.Int        `json:"count"`
	Amount *big.Int        `json:"amount"`
	State  DelegationState `json:"state"`
}

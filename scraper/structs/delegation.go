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
	ETHBlockHeight   uint64          `json:"eth_block_height"`
	Amount           *big.Int        `json:"amount"`
	DelegationPeriod *big.Int        `json:"delegationPeriod"`
	Created          time.Time       `json:"created"`
	Started          *big.Int        `json:"started"`
	Finished         *big.Int        `json:"finished"`
	Info             string          `json:"info"`
	State            DelegationState `json:"state"`
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
		return "unknown"
	}
}

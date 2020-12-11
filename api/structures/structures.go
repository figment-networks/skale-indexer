package structures

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)



type Delegation struct {
	ID               *big.Int        `json:"id"`
	Holder           common.Address  `json:"holder"`
	ValidatorID      *big.Int        `json:"validatorId"`
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

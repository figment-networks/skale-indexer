package structures

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Validator struct {
	ID                      *big.Int       `json:"id"`
	Name                    string         `json:"name"`
	ValidatorAddress        common.Address `json:"validatorAddress"`
	RequestedAddress        common.Address `json:"requestedAddress"`
	Description             string         `json:"description"`
	FeeRate                 *big.Int       `json:"feeRate"`
	RegistrationTime        time.Time      `json:"registrationTime"`
	MinimumDelegationAmount *big.Int       `json:"minimumDelegationAmount"`
	AcceptNewRequests       bool           `json:"acceptNewRequests"`
	Authorized              bool           `json:"authorized"`
}

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

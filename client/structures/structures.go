package structures

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type ContractEvent struct {
	Type         string
	ContractName string
	Address      common.Address
	Height       uint64
	TxHash       common.Hash
	Params       map[string]interface{}
}

type EthEvent struct {
	Address common.Address
}

type Validator struct {
	ID                      *big.Int       `json:"id"`
	Name                    string         `json:"name"`
	ValidatorAddress        common.Address `json:"validatorAddress"`
	RequestedAddress        common.Address `json:"requestedAddress"`
	Description             string         `json:"description"`
	FeeRate                 *big.Int       `json:"feeRate"`
	RegistrationTime        *big.Int       `json:"registrationTime"`
	MinimumDelegationAmount *big.Int       `json:"minimumDelegationAmount"`
	AcceptNewRequests       bool           `json:"acceptNewRequests"`
	Authorized              bool           `json:"authorized"`
}

type NodeStatus uint

const (
	NodeStatusActive NodeStatus = iota
	NodeStatusLeaving
	NodeStatusLeft
	NodeStatusInMaintenance
)

type Node struct {
	ID             *big.Int   `json:"id"`
	Name           string     `json:"name"`
	IP             [4]byte    `json:"ip"`
	PublicIP       [4]byte    `json:"publicIP"`
	Port           uint16     `json:"port"`
	PublicKey      *big.Int   `json:"publicKey"`
	StartBlock     *big.Int   `json:"startBlock"`
	NextRewardDate time.Time  `json:"nextRewardDate"`
	LastRewardDate time.Time  `json:"lastRewardDate"`
	FinishTime     *big.Int   `json:"finishTime"`
	Status         NodeStatus `json:"nodeStatus"`
	ValidatorID    *big.Int   `json:"validatorID"`
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

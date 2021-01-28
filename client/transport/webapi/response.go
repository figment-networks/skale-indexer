package webapi

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// swagger:model
type ContractEvents []ContractEvent

// swagger:model
type Nodes []Node

// swagger:model
type Validators []Validator

// swagger:model
type ValidatorStatistics []ValidatorStatistic

// swagger:model
type Delegations []Delegation

// swagger:model
type Accounts []Account

// ContractEvent a set of fields to show returned events (contract events) by search
// swagger:model
type ContractEvent struct {
	// ID - Identification on database, not the mainnet
	ID uuid.UUID `json:"id"`
	// ContractName - Name of the contract
	//
	// example: delegation_controller
	ContractName string `json:"contract_name"`
	// EventName - name of the event in contract_name
	//
	// example: ValidatorRegistered
	EventName string `json:"event_name"`
	// ContractAddress - Address of the contract (Address represents the 20 byte address of an Ethereum account)
	//
	// package: github.com/ethereum/go-ethereum/common
	ContractAddress common.Address `json:"contract_address"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Time - Event time
	//
	// A Time represents an instant in time with nanosecond precision.
	Time time.Time `json:"time"`
	// TransactionHash - transaction where the event occurred
	//
	// package: github.com/ethereum/go-ethereum/common
	// Hash represents the 32 byte Keccak256 hash of arbitrary data
	// format: string
	TransactionHash common.Hash `json:"transaction_hash"`
	// Removed - indicates whether the event is removed on SKALE
	Removed bool `json:"removed"`
	// Event params
	//
	// bounty id array or address array
	Params map[string]interface{} `json:"params"`
}

// Delegation a set of fields to show returned delegations by search
// swagger:model
type Delegation struct {
	// DelegationID - the index of delegation in SKALE deployed smart contract
	//
	// package: math/big
	DelegationID *big.Int `json:"id"`
	// Holder - Address of the token holder (Address represents the 20 byte address of an Ethereum account.)
	//
	// package: github.com/ethereum/go-ethereum/common
	// format: [20]byte
	Holder common.Address `json:"holder"`
	// TransactionHash - transaction where delegation updated ( Hash represents the 32 byte Keccak256 hash of arbitrary data)
	//
	// package: github.com/ethereum/go-ethereum/common
	// format: [32]byte
	TransactionHash common.Hash `json:"transaction_hash"`
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// package: math/big
	ValidatorID *big.Int `json:"validator_id"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Amount - delegation amount SKL unit
	Amount string `json:"amount"`
	// Period - The duration delegation as chosen by the delegator
	//
	// package: math/big
	Period *big.Int `json:"period"`
	// Created - Creation time at ETH mainnet
	//
	// package: time
	// A Time represents an instant in time with nanosecond precision.
	Created time.Time `json:"created"`
	// Started - started  epoch
	//
	// package: math/big
	Started *big.Int `json:"started"`
	// Finished - finished  epoch
	//
	// package: math/big
	Finished *big.Int `json:"finished"`
	// Info - delegation information
	Info string `json:"info"`
	// State - delegation state
	State string `json:"state"`
}

// Node a set of fields to show returned nodes by search
// swagger:model
type Node struct {
	// NodeID - the index of node in SKALE deployed smart contract
	//
	// package: math/big
	NodeID *big.Int `json:"id"`
	// Name - node name
	Name string `json:"name"`
	// IP - node ip
	IP string `json:"ip"`
	// PublicIP - node public ip
	PublicIP string `json:"public_ip"`
	// Port - node port
	Port uint16 `json:"port"`
	// StartBlock - starting block height on ETH mainnet
	//
	// package: math/big
	StartBlock *big.Int `json:"start_block"`
	// NextRewardDate - next reward time
	//
	// package: time
	NextRewardDate time.Time `json:"next_reward_date"`
	// LastRewardDate - last reward time
	//
	// package: time
	LastRewardDate time.Time `json:"last_reward_date"`
	// FinishTime - finish time
	//
	// package: math/big
	FinishTime *big.Int `json:"finish_time"`
	// ValidatorID - validator Id on SKALE network
	//
	// package: math/big
	ValidatorID *big.Int `json:"validator_id"`
	// Status - node status
	Status string `json:"status"`
}

// Validator a set of fields to show returned validators by search
// swagger:model
type Validator struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// package: math/big
	ValidatorID *big.Int `json:"id"`
	// Name - validator name
	Name string `json:"name"`
	// Description - validator description
	Description string `json:"description"`
	// ValidatorAddress - validator address on SKALE (Address represents the 20 byte address of an Ethereum account)
	//
	// package: github.com/ethereum/go-ethereum/common
	// format: [20]byte
	ValidatorAddress common.Address `json:"validator_address"`
	// RequestedAddress - requested address on SKALE (Address represents the 20 byte address of an Ethereum account)
	//
	// package: github.com/ethereum/go-ethereum/common
	// format: [20]byte
	RequestedAddress common.Address `json:"requested_address"`
	// FeeRate - fee rate
	//
	// package: math/big
	FeeRate *big.Int `json:"fee_rate"`
	// RegistrationTime - registration time to network
	//
	// package: time
	RegistrationTime time.Time `json:"registration_time"`
	// MinimumDelegationAmount - minimum delegation amount i.e. MDR
	//
	// package: math/big
	MinimumDelegationAmount *big.Int `json:"minimum_delegation_amount"`
	// AcceptNewRequests - shows whether validator accepts new requests or not
	AcceptNewRequests bool `json:"accept_new_requests"`
	// Authorized - shows whether validator is authorized or not
	Authorized bool `json:"authorized"`
	// ActiveNodes - number of active nodes attached to the validator
	ActiveNodes uint `json:"active_nodes"`
	// LinkedNodes - number of all nodes attached to the validator
	LinkedNodes uint `json:"linked_nodes"`
	// Staked - total stake amount
	Staked string `json:"staked"`
}

// ValidatorStatistic validator statistic value in given block
// swagger:model
type ValidatorStatistic struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	// package: math/big
	ValidatorID *big.Int `json:"id"`
	// Amount - statistics amount
	Amount string `json:"amount"`
	// BlockHeight - block height on ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Time - block timestamp on ETH mainnet
	Time time.Time `json:"time"`
	// Type - statistics type
	Type string `json:"type"`
}

// Account structure representing ethereum account used in SKALE
// swagger:model
type Account struct {
	// Address - account address (Address represents the 20 byte address of an Ethereum account)
	//
	// package: github.com/ethereum/go-ethereum/common
	// format: [20]byte
	Address common.Address `json:"address"`
	// Type - account type
	Type string `json:"type"`
}

// SystemEvent event information for reporting some activities in chain
// swagger:model
type SystemEvent struct {
	Height      uint64          `json:"height"`
	Time        time.Time       `json:"time"`
	Kind        string          `json:"kind"`
	SenderID    uint64          `json:"sender_id"`
	RecipientID uint64          `json:"recipient_id"`
	Sender      common.Address  `json:"sender"`
	Recipient   common.Address  `json:"recipient"`
	Data        SystemEventData `json:"data"`
}

// SystemEventData value for SystemEvent
// swagger:model
type SystemEventData struct {
	Before big.Int   `json:"before"`
	After  big.Int   `json:"after"`
	Change big.Float `json:"change"`
}

// ApiError a set of fields to show error
// swagger:model
type ApiError struct {
	// Error - error message from api
	Error string `json:"error"`
	// Code - http code
	Code int `json:"code"`
}

func newApiError(err error, code int) []byte {
	resp, _ := json.Marshal(ApiError{
		Error: err.Error(),
		Code:  code,
	})
	return resp
}

package webapi

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// ContractEvent a set of fields to show returned events (contract events) by search
type ContractEvent struct {
	// ID - Identification on database, not the mainnet
	//
	// package: github.com/google/uuid
	// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
	// format: [16]byte
	ID uuid.UUID `json:"id"`
	// ContractName - Name of the contract
	//
	// example: "delegation_controller", "validator_service etc"
	ContractName string `json:"contract_name"`
	// EventName - name of the event in contract_name
	//
	// example: "ValidatorRegistered", "WithdrawBounty"
	EventName string `json:"event_name"`
	// ContractAddress - Address of the contract
	//
	// package: github.com/ethereum/go-ethereum/common
	// Address represents the 20 byte address of an Ethereum account
	// format: [20]byte
	ContractAddress common.Address `json:"contract_address"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Time - Event time
	//
	// package: time
	// A Time represents an instant in time with nanosecond precision.
	Time time.Time `json:"time"`
	// TransactionHash - transaction where the event occurred
	//
	// package: github.com/ethereum/go-ethereum/common
	// Hash represents the 32 byte Keccak256 hash of arbitrary data
	// format: [32]byte
	TransactionHash common.Hash `json:"transaction_hash"`
	// Removed - indicates whether the event is removed on SKALE
	Removed bool `json:"removed"`
	// Event params
	//
	// bounty id array or address array
	Params map[string]interface{} `json:"params"`
}

// Delegation a set of fields to show returned delegations by search
type Delegation struct {
	// DelegationID - the index of delegation in SKALE deployed smart contract
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	DelegationID *big.Int `json:"id"`
	// Holder - Address of the token holder
	//
	// package: github.com/ethereum/go-ethereum/common
	// Address represents the 20 byte address of an Ethereum account
	// format: [20]byte
	Holder common.Address `json:"holder"`
	// TransactionHash - transaction where delegation updated
	//
	// package: github.com/ethereum/go-ethereum/common
	// Hash represents the 32 byte Keccak256 hash of arbitrary data
	// format: [32]byte
	TransactionHash common.Hash `json:"transaction_hash"`
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	ValidatorID *big.Int `json:"validator_id"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Amount - delegation amount SKL unit
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	Amount *big.Int `json:"amount"`
	// Period - The duration delegation as chosen by the delegator
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	Period *big.Int `json:"period"`
	// Created - Creation time at ETH mainnet
	//
	// package: time
	// A Time represents an instant in time with nanosecond precision.
	Created time.Time `json:"created"`
	// Started - started  epoch
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	Started *big.Int `json:"started"`
	// Finished - finished  epoch
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	Finished *big.Int `json:"finished"`
	// Info - delegation information
	Info string `json:"info"`
}

// Node a set of fields to show returned nodes by search
type Node struct {
	// NodeID - the index of node in SKALE deployed smart contract
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
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
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	StartBlock *big.Int `json:"start_block"`
	// NextRewardDate - next reward date
	//
	// package: time
	// A Time represents an instant in time with nanosecond precision.
	NextRewardDate time.Time `json:"next_reward_date"`
	// LastRewardDate - last reward date
	//
	// package: time
	// A Time represents an instant in time with nanosecond precision.
	LastRewardDate time.Time `json:"last_reward_date"`
	// FinishTime - finish time
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	FinishTime *big.Int `json:"finish_time"`
	// ValidatorID - validator Id on SKALE network
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	ValidatorID *big.Int `json:"validator_id"`
}

// Validator a set of fields to show returned validators by search
type Validator struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	ValidatorID *big.Int `json:"id"`
	// Name - validator name
	Name string `json:"name"`
	// Description - validator description
	Description string `json:"description"`
	// ValidatorAddress - validator address on SKALE
	//
	// package: github.com/ethereum/go-ethereum/common
	// Address represents the 20 byte address of an Ethereum account
	// format: [20]byte
	ValidatorAddress common.Address `json:"validator_address"`
	// RequestedAddress - requested address on SKALE
	//
	// package: github.com/ethereum/go-ethereum/common
	// Address represents the 20 byte address of an Ethereum account
	// format: [20]byte
	RequestedAddress common.Address `json:"requested_address"`
	// FeeRate - fee rate
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	FeeRate *big.Int `json:"fee_rate"`
	// RegistrationTime - registration time to network
	//
	// package: time
	// A Time represents an instant in time with nanosecond precision.
	RegistrationTime time.Time `json:"registration_time"`
	// MinimumDelegationAmount - minimum delegation amount i.e. MDR
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
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
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	Staked  *big.Int `json:"staked"`
	Pending *big.Int `json:"pending"`
	Rewards *big.Int `json:"rewards"`
}

// ValidatorStatistic a set of fields to show returned validator statistics by search
type ValidatorStatistic struct {
	// ValidatorID - the index of validator in SKALE deployed smart contract
	//
	// package: math/big
	// An Int represents a signed multi-precision integer. The zero value for an Int represents the value 0.
	ValidatorID *big.Int `json:"id"`
	// Amount - statistics amount
	Amount string `json:"amount"`
	// BlockHeight - starting block height on ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Type - statistics type
	Type string `json:"type"`
}

// Account a set of fields to show returned accounts by search
type Account struct {
	// Address - account address
	//
	// package: github.com/ethereum/go-ethereum/common
	// Address represents the 20 byte address of an Ethereum account
	// format: [20]byte
	Address common.Address `json:"address"`
	// Type - account type
	Type string `json:"type"`
}

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
	ContractAddress common.Address `json:"contract_address"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Time - Event time
	Time time.Time `json:"time"`
	// TransactionHash - transaction where the event occurred
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
	// DelegationID - delegation Id on SKALE network
	DelegationID *big.Int `json:"id"`
	// Holder - Address of the token holder
	Holder common.Address `json:"holder"`
	// TransactionHash - transaction where delegation updated
	TransactionHash common.Hash `json:"transaction_hash"`
	// ValidatorID - validator Id on SKALE network
	ValidatorID *big.Int `json:"validator_id"`
	// BlockHeight - Block number at ETH mainnet
	BlockHeight uint64 `json:"block_height"`
	// Amount - delegation amount SKL unit
	Amount *big.Int `json:"amount"`
	// Period - The duration delegation as chosen by the delegator
	Period *big.Int `json:"period"`
	// Created - Creation time at ETH mainnet
	Created time.Time `json:"created"`
	// Started - started  epoch
	Started *big.Int `json:"started"`
	// Finished - finished  epoch
	Finished *big.Int `json:"finished"`
	// Info - delegation information
	Info string `json:"info"`
}

// Node a set of fields to show returned nodes by search
type Node struct {
	// NodeID - node Id on SKALE network
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
	StartBlock *big.Int `json:"start_block"`
	// NextRewardDate - next reward date
	NextRewardDate time.Time `json:"next_reward_date"`
	// LastRewardDate - last reward date
	LastRewardDate time.Time `json:"last_reward_date"`
	// FinishTime - finish time
	FinishTime *big.Int `json:"finish_time"`
	// ValidatorID - validator Id on SKALE network
	ValidatorID *big.Int `json:"validator_id"`
}

// Validator a set of fields to show returned validators by search
type Validator struct {
	// ValidatorID - validator Id on SKALE network
	ValidatorID *big.Int `json:"id"`
	// Name - validator name
	Name string `json:"name"`
	// Description - validator description
	Description string `json:"description"`
	// ValidatorAddress - validator address on SKALE
	ValidatorAddress common.Address `json:"validator_address"`
	// RequestedAddress - requested address on SKALE
	RequestedAddress common.Address `json:"requested_address"`
	// FeeRate - fee rate
	FeeRate *big.Int `json:"fee_rate"`
	// RegistrationTime - registration time to network
	RegistrationTime time.Time `json:"registration_time"`
	// MinimumDelegationAmount - minimum delegation amount i.e. MDR
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
	Staked  *big.Int `json:"staked"`
	Pending *big.Int `json:"pending"`
	Rewards *big.Int `json:"rewards"`
}

// ValidatorStatistic a set of fields to show returned validator statistics by search
type ValidatorStatistic struct {
	// ValidatorID - validator Id on SKALE network
	ValidatorID uint64 `json:"id"`
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
	Address common.Address `json:"address"`
	// Type - account type
	Type string `json:"type"`
}

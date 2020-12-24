package webapi

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ContractEventAPI struct {
	ID              uuid.UUID              `json:"id"`
	ContractName    string                 `json:"contract_name"`
	EventName       string                 `json:"event_name"`
	ContractAddress common.Address         `json:"contract_address"`
	BlockHeight     uint64                 `json:"block_height"`
	Time            time.Time              `json:"time"`
	TransactionHash common.Hash            `json:"transaction_hash"`
	Removed         bool                   `json:"removed"`
	Params          map[string]interface{} `json:"params"`
}

type DelegationAPI struct {
	DelegationID     *big.Int       `json:"delegation_id"`
	Holder           common.Address `json:"holder"`
	TransactionHash  common.Hash    `json:"transaction_hash"`
	ValidatorID      *big.Int       `json:"validator_id"`
	BlockHeight      uint64         `json:"block_height"`
	Amount           *big.Int       `json:"amount"`
	DelegationPeriod *big.Int       `json:"delegation_period"`
	Created          time.Time      `json:"created"`
	Started          *big.Int       `json:"started"`
	Finished         *big.Int       `json:"finished"`
	Info             string         `json:"info"`
}

type NodeAPI struct {
	NodeID         *big.Int  `json:"node_id"`
	Name           string    `json:"name"`
	IP             string    `json:"ip"`
	PublicIP       string    `json:"public_ip"`
	Port           uint16    `json:"port"`
	StartBlock     *big.Int  `json:"start_block"`
	NextRewardDate time.Time `json:"next_reward_date"`
	LastRewardDate time.Time `json:"last_reward_date"`
	FinishTime     *big.Int  `json:"finish_time"`
	ValidatorID    *big.Int  `json:"validator_id"`
}

// TODO: change api response with this
type ValidatorAPI struct {
	ValidatorID             *big.Int       `json:"validator_id"`
	Name                    string         `json:"name"`
	ValidatorAddress        common.Address `json:"validator_address"`
	RequestedAddress        common.Address `json:"requested_address"`
	Description             string         `json:"description"`
	FeeRate                 *big.Int       `json:"fee_rate"`
	RegistrationTime        time.Time      `json:"registration_time"`
	MinimumDelegationAmount *big.Int       `json:"minimum_delegation_amount"`
	AcceptNewRequests       bool           `json:"accept_new_requests"`
	Authorized              bool           `json:"authorized"`
	Active                  bool           `json:"active"`
	ActiveNodes             int            `json:"active_nodes"`
	LinkedNodes             int            `json:"linked_nodes"`
	Staked                  uint64         `json:"staked"`
	Pending                 uint64         `json:"pending"`
	Rewards                 uint64         `json:"rewards"`
}

// TODO: change api response with this
type ValidatorStatisticsAPI struct {
	ValidatorId uint64 `json:"validator_id"`
	Amount      uint64 `json:"amount"`
	BlockHeight uint64 `json:"block_height"`
}

type AccountAPI struct {
	Address     common.Address `json:"address"`
	AccountType string         `json:"account_type"`
}

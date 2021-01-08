package webapi

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ContractEvent struct {
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

type Delegation struct {
	DelegationID    *big.Int       `json:"id"`
	Holder          common.Address `json:"holder"`
	TransactionHash common.Hash    `json:"transaction_hash"`
	ValidatorID     *big.Int       `json:"validator_id"`
	BlockHeight     uint64         `json:"block_height"`
	Amount          *big.Int       `json:"amount"`
	Period          *big.Int       `json:"period"`
	Created         time.Time      `json:"created"`
	Started         *big.Int       `json:"started"`
	Finished        *big.Int       `json:"finished"`
	Info            string         `json:"info"`
}

type Node struct {
	NodeID         *big.Int  `json:"id"`
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

type Validator struct {
	ValidatorID             *big.Int       `json:"id"`
	Name                    string         `json:"name"`
	Description             string         `json:"description"`
	ValidatorAddress        common.Address `json:"validator_address"`
	RequestedAddress        common.Address `json:"requested_address"`
	FeeRate                 *big.Int       `json:"fee_rate"`
	RegistrationTime        time.Time      `json:"registration_time"`
	MinimumDelegationAmount *big.Int       `json:"minimum_delegation_amount"`
	AcceptNewRequests       bool           `json:"accept_new_requests"`
	Authorized              bool           `json:"authorized"`
	ActiveNodes             uint           `json:"active_nodes"`
	LinkedNodes             uint           `json:"linked_nodes"`
	Staked                  *big.Int       `json:"staked"`
	Pending                 *big.Int       `json:"pending"`
	Rewards                 *big.Int       `json:"rewards"`
}

type ValidatorStatistic struct {
	ValidatorID uint64 `json:"id"`
	Amount      string `json:"amount"`
	BlockHeight uint64 `json:"block_height"`
	Type        string `json:"type"`
}

type Account struct {
	Address common.Address `json:"address"`
	Type    string         `json:"type"`
}

package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type Validator struct {
	ID                      string         `json:"id"`
	CreatedAt               time.Time      `json:"created_at"`
	// TODO: remove updated_at
	UpdatedAt               time.Time      `json:"updated_at"`
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
	ETHBlockHeight          uint64         `json:"eth_block_height"`
}

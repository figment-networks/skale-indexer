package structs

import (
	"time"
)

type Validator struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Address     []int     `json:"address"`
	Description string    `json:"description"`
	FeeRate     uint64    `json:"fee_rate"`
	Active      bool      `json:"active"`
	ActiveNodes int       `json:"active_nodes"`
	Staked      uint64    `json:"staked"`
	Pending     uint64    `json:"pending"`
	Rewards     uint64    `json:"rewards"`
	Data        []Data 	  `json:"data"`
}

type Data struct {
	RequestedAddress        string    `json:"requested_address"`
	RegistrationTime        time.Time `json:"registration_time"`
	MinimumDelegationAmount uint64    `json:"minimum_delegation_amount"`
	AcceptNewRequests       bool      `json:"accept_new_requests"`
	Trusted                 bool      `json:"trusted"`
}

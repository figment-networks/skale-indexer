package structs

import (
	"time"
)

type Validator struct {
	ID                      string    `json:"id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	Name                    string    `json:"name"`
	ValidatorAddress        string    `json:"validator_address"`
	RequestedAddress        string    `json:"requested_address"`
	Description             string    `json:"description"`
	FeeRate                 uint64    `json:"fee_rate"`
	RegistrationTime        time.Time `json:"registration_time"`
	MinimumDelegationAmount uint64    `json:"minimum_delegation_amount"`
	AcceptNewRequests       bool      `json:"accept_new_requests"`
}

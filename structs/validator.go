package structs

import "../types"

type Validator struct {
	ID                      *types.ID   `json:"id"`
	CreatedAt               *types.Time `json:"created_at"`
	UpdatedAt               *types.Time `json:"updated_at"`
	Name                    *string     `json:"name"`
	ValidatorAddress        *string     `json:"validator_address"`
	RequestedAddress        *string     `json:"requested_address"`
	Description             *string     `json:"description"`
	FeeRate                 *uint64     `json:"fee_rate"`
	RegistrationTime        *uint64     `json:"registration_time"`
	MinimumDelegationAmount *uint64     `json:"minimum_delegation_amount"`
	AcceptNewRequests       *uint64     `json:"accept_new_requests"`
}

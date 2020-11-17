package structs

import (
	"time"
)

type Delegation struct {
	ID               *string    `json:"id"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Holder           *string    `json:"holder"`
	ValidatorId      *uint64    `json:"validator_id"`
	Amount           *uint64    `json:"amount"`
	DelegationPeriod *uint64    `json:"delegation_period"`
	Created          *time.Time `json:"created"`
	Started          *time.Time `json:"started"`
	Finished         *time.Time `json:"finished"`
	Info             *string    `json:"info"`
}

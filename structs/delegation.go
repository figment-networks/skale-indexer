package structs

import "../types"

type Delegation struct {
	ID               types.ID   `json:"id"`
	CreatedAt        types.Time `json:"created_at"`
	UpdatedAt        types.Time `json:"updated_at"`
	Holder           string     `json:"holder"`
	ValidatorId      uint64     `json:"validator_id"`
	Amount           uint64     `json:"amount"`
	DelegationPeriod uint64     `json:"delegation_period"`
	Created          uint64     `json:"created"`
	Started          uint64     `json:"started"`
	Finished         uint64     `json:"finished"`
	Info             uint64     `json:"info"`
}

package structs

import "github.com/figment-networks/skale-indexer/types"

type Delegation struct {
	ID               types.ID   `json:"id"`
	Holder           string     `json:"holder"`
	ValidatorId      uint64     `json:"validator_id"`
	Amount           uint64     `json:"amount"`
	DelegationPeriod uint64     `json:"delegation_period"`
	Created          uint64     `json:"created"`
	Started          uint64     `json:"started"`
	Finished         uint64     `json:"finished"`
	Info             uint64     `json:"info"`
	CreatedAt        types.Time `json:"created_at"`
	UpdatedAt        types.Time `json:"updated_at"`
}

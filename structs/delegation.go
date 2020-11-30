package structs

import (
	"time"
)

type Delegation struct {
	ID                   string           `json:"id"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	Holder               uint64           `json:"holder"`
	ValidatorId          uint64           `json:"validator_id"`
	Amount               uint64           `json:"amount"`
	DelegationPeriod     uint64           `json:"delegation_period"`
	Created              time.Time        `json:"created"`
	Started              time.Time        `json:"started"`
	Finished             time.Time        `json:"finished"`
	Info                 string           `json:"info"`
	Status               DelegationStatus `json:"status"`
	SmartContractIndex   uint64           `json:"smart_contract_index"`
	SmartContractAddress uint64           `json:"smart_contract_address"`
}

type DelegationStatus int

const (
	Proposed DelegationStatus = iota + 1
	Accepted
	Canceled
	Rejected
	Delegated
	UndelegatedRequested
	Completed
	Pending     // not available in the source code
	UnDelegated // not available in the source code
)

func (k DelegationStatus) String() string {
	switch k {
	case Proposed:
		return "PROPOSED"
	case Accepted:
		return "ACCEPTED"
	case Canceled:
		return "CANCELED"
	case Rejected:
		return "REJECTED"
	case Delegated:
		return "DELEGATED"
	case UndelegatedRequested:
		return "UNDELEGATION_REQUESTED"
	case Completed:
		return "COMPLETED"
	case Pending:
		return "PENDING"
	case UnDelegated:
		return "UNDELEGATED"
	default:
		return "unknown"
	}
}

package structs

import (
	"time"
)

type DelegationStateStatistics struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      uint      `json:"status"`
	ValidatorId uint64    `json:"validator_id"`
	Amount      uint64    `json:"amount"`
}

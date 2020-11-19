package structs

import (
	"time"
)

type DelegationEvent struct {
	ID           *string    `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DelegationId *string    `json:"delegation_id"`
	EventName    *string    `json:"event_name"`
	EventTime    *time.Time `json:"event_time"`
}

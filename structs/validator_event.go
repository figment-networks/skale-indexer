package structs

import (
	"time"
)

type ValidatorEvent struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ValidatorId string    `json:"validator_id"`
	EventName   string    `json:"event_name"`
	EventTime   time.Time `json:"event_time"`
}

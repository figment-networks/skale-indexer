package structs

import (
	"time"
)

type Event struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	BlockHeight          int64     `json:"block_height"`
	SmartContractAddress string    `json:"smart_contract_address"`
	TransactionIndex     int64     `json:"transaction_index"`
	EventType            string    `json:"event_type"`
	EventName            string    `json:"event_name"`
	EventTime            time.Time `json:"event_time"`
}

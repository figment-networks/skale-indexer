package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	EventInfo            EventInfo `json:"event_info"`
}

type EventInfo struct {
	Wallet      string    `json:"wallet"`
	Holder      string    `json:"holder"`
	Destination []Address `json:"destination"`
	ValidatorId uint64    `json:"validator_id"`
	Amount      uint64    `json:"amount"`
}

func (a EventInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *EventInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

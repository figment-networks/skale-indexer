package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Validator struct {
	ID           string       `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Name         string       `json:"name"`
	ValidatorId  uint64       `json:"validator_id"`
	Address      []Address    `json:"address"`
	Description  string       `json:"description"`
	FeeRate      uint64       `json:"fee_rate"`
	Active       bool         `json:"active"`
	ActiveNodes  int          `json:"active_nodes"`
	LinkedNodes  int          `json:"linked_nodes"`
	Staked       uint64       `json:"staked"`
	Pending      uint64       `json:"pending"`
	Rewards      uint64       `json:"rewards"`
	OptionalInfo OptionalInfo `json:"optional_info"`
}

type OptionalInfo struct {
	Data []Data `json:"data"`
}

type Data struct {
	RequestedAddress        string    `json:"requested_address"`
	RegistrationTime        time.Time `json:"registration_time"`
	MinimumDelegationAmount uint64    `json:"minimum_delegation_amount"`
	AcceptNewRequests       bool      `json:"accept_new_requests"`
	Trusted                 bool      `json:"trusted"`
}

func (a OptionalInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *OptionalInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type Address uint64

func (a *Address) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

package structs

import (
	"time"
)

type Validator struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Description string    `json:"description"`
	FeeRate     uint64    `json:"fee_rate"`
	Active      bool      `json:"active"`
	ActiveNodes int       `json:"active_nodes"`
	Staked      uint64    `json:"staked"`
	Pending     uint64    `json:"pending"`
	Rewards     uint64    `json:"rewards"`
	//Data        Data 	  `json:"data"`
}

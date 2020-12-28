package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	Address     common.Address `json:"address"`
	AccountType AccountType    `json:"account_type"`
}

type AccountType uint

const (
	AccountTypeDefault AccountType = iota
	AccountTypeDelegator
	AccountTypeValidator
)

func (k AccountType) String() string {
	switch k {
	case AccountTypeDefault:
		return "default"
	case AccountTypeDelegator:
		return "delegator"
	case AccountTypeValidator:
		return "validator"
	default:
		return "unknown"
	}
}

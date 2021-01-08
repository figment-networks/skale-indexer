package structs

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type Account struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	Address     common.Address `json:"address"`
	AccountType AccountType    `json:"account_type"`
}

type AccountType string

const (
	AccountTypeDefault AccountType = "default"
	AccountTypeDelegator AccountType= "delegator"
	AccountTypeValidator AccountType = "validator"
)
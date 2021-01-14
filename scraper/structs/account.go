package structs

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	Address   common.Address `json:"address"`
	Type      AccountType    `json:"type"`
}

type AccountType string

const (
	AccountTypeDefault   AccountType = "default"
	AccountTypeDelegator AccountType = "delegator"
	AccountTypeValidator AccountType = "validator"
)

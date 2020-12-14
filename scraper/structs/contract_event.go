package structs

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type ContractEvent struct {
	ID              string                 `json:"id"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	ContractName    string                 `json:"contract_name"`
	EventName       string                 `json:"event_name"`
	ContractAddress common.Address         `json:"contract_address"`
	BlockHeight     uint64                 `json:"block_height"`
	Time            time.Time              `json:"time"`
	TransactionHash common.Hash            `json:"transaction_hash"`
	Removed         bool                   `json:"removed"`
	Params          map[string]interface{} `json:"params"`
	BoundType       string                 `json:"boundType"`
	BoundAddress    common.Address         `json:"boundAddress"`
}

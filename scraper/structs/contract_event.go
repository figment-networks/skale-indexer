package structs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ContractEvent struct {
	ID              uuid.UUID              `json:"id"`
	ContractName    string                 `json:"contract_name"`
	EventName       string                 `json:"event_name"`
	ContractAddress common.Address         `json:"contract_address"`
	BlockHeight     uint64                 `json:"block_height"`
	Time            time.Time              `json:"time"`
	TransactionHash common.Hash            `json:"transaction_hash"`
	Removed         bool                   `json:"removed"`
	Params          map[string]interface{} `json:"params"`
	BoundType       string                 `json:"bound_type"`
	BoundID         []big.Int              `json:"bound_id"`
	BoundAddress    []common.Address       `json:"bound_address"`
}

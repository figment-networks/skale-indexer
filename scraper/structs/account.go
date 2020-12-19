package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"math/big"
	"time"
)

type Account struct {
	ID        		uuid.UUID      `json:"id"`
	CreatedAt 		time.Time      `json:"created_at"`
	Address   		common.Address `json:"address"`
	BoundKind 		string         `json:"bound_kind"`
	BoundID   		*big.Int       `json:"bound_id"`
	BlockHeight     uint64          `json:"block_height"`
}

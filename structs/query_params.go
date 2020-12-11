package structs

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type QueryParams struct {
	Id               string
	ValidatorId      uint64
	ValidatorAddress common.Address
	TimeFrom         time.Time
	TimeTo           time.Time
}

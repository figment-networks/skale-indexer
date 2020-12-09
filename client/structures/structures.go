package structures

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type ContractEvent struct {
	Type         string
	ContractName string
	Time         time.Time
	Address      common.Address
	Height       uint64
	TxHash       common.Hash
	Removed      bool
	Params       map[string]interface{}
}

type EthEvent struct {
	Address common.Address
}

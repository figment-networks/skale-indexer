package transport

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type EthereumTransport interface {
	Dial(ctx context.Context) (err error)
	Close(ctx context.Context)
	GetLogs(ctx context.Context, from, to big.Int) (logs []types.Log, err error)
}

package transport

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthereumTransport interface {
	Dial(ctx context.Context) (err error)
	Close(ctx context.Context)
	GetLogs(ctx context.Context, from, to big.Int, contracts []common.Address) (logs []types.Log, err error)
	GetBlockHeader(ctx context.Context, height *big.Int) (h *types.Header, err error)
	GetBoundContractCaller(ctx context.Context, address common.Address, a abi.ABI) *bind.BoundContract
}

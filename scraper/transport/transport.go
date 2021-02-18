package transport

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var ErrEmptyResponse = errors.New("Returned Empty Response (Reverted 0x)")

type BoundContractCaller interface {
	GetContract() *bind.BoundContract
	RawCall(ctx context.Context, opts *bind.CallOpts, method string, params ...interface{}) (output []byte, err error)
	AbiUnpack(method string, data []byte) (res []interface{}, err error)
}

type EthereumTransport interface {
	Dial(ctx context.Context) (err error)
	Close(ctx context.Context)
	GetLogs(ctx context.Context, from, to big.Int, contracts []common.Address) (logs []types.Log, err error)
	GetBlockHeader(ctx context.Context, height *big.Int) (h *types.Header, err error)
	GetBoundContractCaller(ctx context.Context, address common.Address, a abi.ABI) BoundContractCaller
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
}

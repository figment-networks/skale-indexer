package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthTransport struct {
	C   *ethclient.Client
	Url string
}

func NewEthTransport(url string) *EthTransport {
	return &EthTransport{Url: url}
}

func (et *EthTransport) Dial(ctx context.Context) (err error) {
	et.C, err = ethclient.DialContext(ctx, et.Url)
	return err
}

func (et *EthTransport) Close(ctx context.Context) {
	et.C.Close()
	return
}

func (et *EthTransport) GetBoundContractCaller(ctx context.Context, address common.Address, a abi.ABI) *bind.BoundContract {
	return bind.NewBoundContract(address, a, et.C, nil, nil)
}

// TODO: validate from-to range
func (et *EthTransport) GetLogs(ctx context.Context, from, to big.Int, contracts []common.Address) (logs []types.Log, err error) {
	fq := ethereum.FilterQuery{
		FromBlock: &from,
		ToBlock:   &to,
	}

	if contracts != nil {
		for _, k := range contracts {
			fq.Addresses = append(fq.Addresses, k)
		}
	}

	return et.C.FilterLogs(ctx, fq)
}

func (et *EthTransport) GetBlockHeader(ctx context.Context, height *big.Int) (h *types.Header, err error) {
	h, err = et.C.HeaderByNumber(ctx, height)
	return h, err
}

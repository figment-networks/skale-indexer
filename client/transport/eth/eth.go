package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
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

package eth

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/figment-networks/skale-indexer/scraper/transport"
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

func (et *EthTransport) GetBoundContractCaller(ctx context.Context, address common.Address, a abi.ABI) transport.BoundContractCaller {
	return &BoundContractC{
		address: address,
		abi:     a,
		ET:      et}

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

func (et *EthTransport) GetCurrentBlockHeight(ctx context.Context) (uint64, error) {
	blockNumber, err := et.C.BlockNumber(ctx)
	return blockNumber, err
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type BoundContractC struct {
	address common.Address
	abi     abi.ABI
	ET      *EthTransport
}

func (bcc *BoundContractC) GetContract() *bind.BoundContract {
	return bind.NewBoundContract(bcc.address, bcc.abi, bcc.ET.C, nil, nil)
}

func (bcc *BoundContractC) RawCall(ctx context.Context, opts *bind.CallOpts, method string, params ...interface{}) (output []byte, err error) {

	// Pack the input, call and unpack the results
	input, err := bcc.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	var (
		msg  = ethereum.CallMsg{From: opts.From, To: &bcc.address, Data: input}
		code []byte
	)

	output, err = bcc.ET.C.CallContract(ctx, msg, opts.BlockNumber)
	if err == nil && len(output) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err = bcc.ET.C.CodeAt(ctx, bcc.address, opts.BlockNumber); err != nil {
			return nil, err
		} else if len(code) == 0 {
			return nil, bind.ErrNoCode
		}
	}

	if err != nil {
		b, err2 := json.Marshal(err)
		if err2 == nil {
			a := map[string]interface{}{}

			err2 := json.Unmarshal(b, &a)
			if err2 == nil {
				if d, ok := a["data"]; ok && d == "Reverted 0x" {
					return output, transport.ErrEmptyResponse
				}

			}
		}
	}
	return output, err
}

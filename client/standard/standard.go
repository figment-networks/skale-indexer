package standard

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
)

func DecodeERC20Events(ctx context.Context, ce structs.ContractEvent) (ceOut structs.ContractEvent, ev error) {
	switch ce.ContractName {
	case "Transfer":
		// Transfer(from, to, value)
		fromI, ok := ce.Params["from"]
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer")
		}
		from, ok := fromI.(common.Address)
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer (from is not Address)")
		}

		toI, ok := ce.Params["to"]
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer")
		}
		to, ok := toI.(common.Address)
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer (to is not Address)")
		}

		ce.BoundAddress = []common.Address{from, to}
	case "Approval":
		// Approval(owner, spender, value)
		ownerI, ok := ce.Params["owner"]
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer")
		}
		owner, ok := ownerI.(common.Address)
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer (owner is not Address)")
		}

		spenderI, ok := ce.Params["spender"]
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer")
		}
		spender, ok := spenderI.(common.Address)
		if !ok {
			return ce, errors.New("structure is not a ERC20 Transfer (spender is not Address)")
		}

		ce.BoundAddress = []common.Address{owner, spender}
	}

	return ce, nil
}

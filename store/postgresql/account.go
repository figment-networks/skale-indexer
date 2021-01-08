package postgresql

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"math/big"
	"strconv"
	"strings"
)

// SaveAccount saves account
func (d *Driver) SaveAccount(ctx context.Context, a structs.Account) error {
	_, err := d.db.Exec(`INSERT INTO accounts ("address", "account_type") 
			VALUES ($1, $2) 
			ON CONFLICT (address)
			DO UPDATE SET
			account_type = EXCLUDED.account_type 
			WHERE accounts.account_type < EXCLUDED.account_type `,
		a.Address.Hash().Big().String(),
		a.AccountType)
	return err
}

// GetAccounts gets accounts
func (d *Driver) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {
	q := `SELECT id, created_at, address, account_type FROM accounts `
	var (
		args   []interface{}
		wherec []string
		i      = 1
	)

	if params.Address != "" {
		wherec = append(wherec, ` address =  $`+strconv.Itoa(i))
		args = append(args, common.HexToAddress(params.Address).Hash().Big().String())
		i++
	}
	if params.Type != "" {
		wherec = append(wherec, ` account_type =  $`+strconv.Itoa(i))
		args = append(args, params.Type)
		i++
	}
	if len(args) > 0 {
		q += ` WHERE `
	}
	q += strings.Join(wherec, " AND ")
	q += ` ORDER BY created_at DESC`

	rows, err := d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		a := structs.Account{}
		var addr []byte
		if err = rows.Scan(&a.ID, &a.CreatedAt, &addr, &a.AccountType); err != nil {
			return nil, err
		}
		p := new(big.Int)
		p.SetString(string(addr), 10)
		a.Address.SetBytes(p.Bytes())

		accounts = append(accounts, a)
	}

	return accounts, nil
}

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"math/big"
)

// SaveAccount saves account
func (d *Driver) SaveAccount(ctx context.Context, a structs.Account) error {
	_, err := d.db.Exec(`INSERT INTO accounts ("address", "account_type") 
			VALUES ($1, $2) 
			ON CONFLICT (address)
			DO UPDATE SET
			account_type = EXCLUDED.account_type `,
		a.Address.Hash().Big().String(),
		a.AccountType)
	return err
}


// GetAccounts gets accounts
func (d *Driver) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {

	q := `SELECT id, created_at, address, account_type FROM accounts `

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q)
	if params.Address != ""{
		q = fmt.Sprintf("%s%s", q, "WHERE address =  $1 ")
		rows, err = d.db.QueryContext(ctx, q, common.HexToAddress(params.Address).Hash().Big().String())
	} else if params.Type != "" {
		q = fmt.Sprintf("%s%s", q, "WHERE account_type =  $1 ")
		rows, err = d.db.QueryContext(ctx, q, params.Type)
	} else {
		rows, err = d.db.QueryContext(ctx, q)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		a := structs.Account{}
		var addr []byte
		if err = rows.Scan(&a.ID, &a.CreatedAt, &addr, &a.AccountType ); err != nil {
			return nil, err
		}
		p := new(big.Int)
		p.SetString(string(addr), 10)
		a.Address.SetBytes(p.Bytes())

		accounts = append(accounts, a)
	}

	return accounts, nil
}

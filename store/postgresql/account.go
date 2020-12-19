package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"math/big"
)

// SaveAccount saves account
func (d *Driver) SaveAccount(ctx context.Context, a structs.Account) error {
	_, err := d.db.Exec(`INSERT INTO accounts ("address", "bound_kind", "bound_id", "block_height") VALUES ($1, $2, $3, $4) `,
		a.Address.Hash().Big().String(),
		a.BoundKind,
		a.BoundID.String(),
		a.BlockHeight)
	return err
}


// GetAccounts gets accounts
func (d *Driver) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {

	q := `SELECT id, created_at, address, bound_type, bound_id, block_height FROM accounts`

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q)


	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		a := structs.Account{}
		var addr []byte
		var bndId uint64
		if err = rows.Scan(&a.ID, &addr, &a.BoundKind, &bndId, &a.BlockHeight ); err != nil {
			return nil, err
		}
		p := new(big.Int)
		p.SetString(string(addr), 10)
		a.Address.SetBytes(p.Bytes())
		a.BoundID = new(big.Int).SetUint64(bndId)

		accounts = append(accounts, a)
	}

	return accounts, nil
}

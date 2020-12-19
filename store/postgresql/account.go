package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"math/big"
)
const (
	orderByBlockHeightA = `ORDER BY block_height  DESC `
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

	q := `SELECT id, created_at, address, bound_kind, bound_id, block_height FROM accounts `

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q)
    // TODO: handle recent param
	if params.Kind != "" && params.Id != ""{
		q = fmt.Sprintf("%s%s%s", q, "WHERE bound_kind =  $1 AND bound_id =  $2 ", orderByBlockHeightA)
		rows, err = d.db.QueryContext(ctx, q, params.Kind, params.Id)
	} else if params.Kind != "" && params.Id == "" {
		q = fmt.Sprintf("%s%s%s", q, "WHERE bound_kind =  $1 ", orderByBlockHeightA)
		rows, err = d.db.QueryContext(ctx, q, params.Kind)
	} else {
		q = fmt.Sprintf("%s%s", q, orderByBlockHeightA)
		rows, err = d.db.QueryContext(ctx, q)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		a := structs.Account{}
		var addr []byte
		var bndId uint64
		if err = rows.Scan(&a.ID, &a.CreatedAt, &addr, &a.BoundKind, &bndId, &a.BlockHeight ); err != nil {
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

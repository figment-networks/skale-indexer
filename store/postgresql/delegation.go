package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// TODO: run explain analyze to check full scan and add required indexes
const (
	// TODO: add started, finished
	getByStatementD = `SELECT id, created_at, delegation_id, holder, validator_id, eth_block_height, amount, delegation_period, created, info, state FROM delegations `
	byDelegationIdD = `WHERE delegation_id =  $1 `
	byCreatedRangeD = `WHERE created between $1 and $2 `
	byValidatorIdD  = `AND validator_id =  $3 `
	orderByCreatedD = `ORDER BY created DESC `

	// for recent
	// TODO: add started, finished
	byRecentEthBlockHeightD = `SELECT  DISTINCT ON (delegation_id) id, created_at, delegation_id, holder, validator_id, eth_block_height, amount, delegation_period, created, info, state
									FROM delegations `
	byRecentValidatorIdD = `WHERE validator_id = $1 `
	orderRecentD         = `ORDER BY delegation_id, eth_block_height DESC`
)

// SaveDelegation saves delegation
func (d *Driver) SaveDelegation(ctx context.Context, dl structs.Delegation) error {
	_, err := d.db.Exec(`INSERT INTO delegations (
				"delegation_id",
				"holder",
				"validator_id",
				"eth_block_height",
				"amount",
				"delegation_period",
				"created",
				"started",
				"finished",
				"info",
				"state")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) `,
		dl.DelegationID.String(),
		dl.Holder.Hash().Big().String(),
		dl.ValidatorID.String(),
		dl.ETHBlockHeight,
		dl.Amount.String(),
		dl.DelegationPeriod.String(),
		dl.Created,
		dl.Started.String(),
		dl.Finished.String(),
		dl.Info,
		dl.State)
	return err
}

// GetDelegations gets delegations by params
func (d *Driver) GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	var q string
	var rows *sql.Rows

	if params.DelegationId != "" {
		q = fmt.Sprintf("%s%s%s", getByStatementD, byDelegationIdD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.DelegationId)
	} else if !params.Recent {
		q = fmt.Sprintf("%s%s", getByStatementD, byCreatedRangeD)
		if params.ValidatorId != "" {
			q = fmt.Sprintf("%s%s%s", q, byValidatorIdD, orderByCreatedD)
			rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, params.ValidatorId)
		} else {
			rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
		}
	} else if params.Recent {
		q = byRecentEthBlockHeightD
		if params.ValidatorId != "" {
			q = fmt.Sprintf("%s%s%s", q, byRecentValidatorIdD, orderRecentD)
			rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
		} else {
			q = fmt.Sprintf("%s%s", q, orderRecentD)
			rows, err = d.db.QueryContext(ctx, q)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.Delegation{}
		var dlgId uint64
		var holder []byte
		var vldId uint64
		var amount []byte
		var dlgPeriod uint64
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlgId, &holder, &vldId, &dlg.ETHBlockHeight, &amount, &dlgPeriod, &dlg.Created, &dlg.Info, &dlg.State)
		if err != nil {
			return nil, err
		}
		dlg.DelegationID = new(big.Int).SetUint64(dlgId)
		h := new(big.Int)
		h.SetString(string(holder), 10)
		dlg.Holder.SetBytes(h.Bytes())
		dlg.ValidatorID = new(big.Int).SetUint64(vldId)
		a := new(big.Int)
		a.SetString(string(amount), 10)
		dlg.Amount = a
		dlg.DelegationPeriod = new(big.Int).SetUint64(dlgPeriod)
		delegations = append(delegations, dlg)
	}
	return delegations, nil
}

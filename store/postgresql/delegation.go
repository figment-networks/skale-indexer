package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

const (
	insertStatementD        = `INSERT INTO delegations ("delegation_id", "holder", "validator_id", "eth_block_height", "amount", "delegation_period", "created", "started",  "finished", "info", "state" ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) `
	getByStatementD         = `SELECT d.id, d.created_at, d.updated_at, d.delegation_id, d.holder, d.validator_id, d.eth_block_height, d.amount, d.delegation_period, d.created, d.started, d.finished, d.info, d.state FROM delegations d `
	byIdD                   = `WHERE d.id =  $1 `
	byValidatorIdD          = `WHERE d.validator_id =  $1 `
	byDateRangeD            = `AND d.created between $1 and $2 `
	byRecentEthBlockHeightD  = `SELECT  DISTINCT ON (delegation_id) d.id, d.created_at, d.updated_at, d.delegation_id, d.holder, d.validator_id, d.eth_block_height, d.amount, d.delegation_period, d.created, d.started, d.finished, d.info, d.state  
									FROM delegations d
								WHERE d.validator_id = $1 AND d.eth_block_height <=$2 
									ORDER BY d.delegation_id, d.eth_block_height DESC`
	orderByCreatedD         = `ORDER BY d.created DESC `
)

// SaveDelegation saves delegation
func (d *Driver) SaveDelegation(ctx context.Context, dl structs.Delegation) error {
	_, err := d.db.Exec(insertStatementD, dl.DelegationID, dl.Holder, dl.ValidatorID, dl.ETHBlockHeight, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info, dl.State)
	return err
}

// GetDelegations gets delegations by params
func (d *Driver) GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error) {
	var q string
	var rows *sql.Rows
	if params.ValidatorId != 0 && !params.Recent {
		q = fmt.Sprintf("%s%s%s", getByStatementD, byValidatorIdD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else if params.ValidatorId != 0 && params.Recent {
		q = fmt.Sprintf("%s%s", byRecentEthBlockHeightD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId, params.ETHBlockHeight)
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		q = fmt.Sprintf("%s%s%s", getByStatementD, byDateRangeD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	} else {
		//TODO: remove by id from api
		q = fmt.Sprintf("%s%s", getByStatementD, byIdD)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.Delegation{}
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.DelegationID, &dlg.Holder, &dlg.ValidatorID, &dlg.ETHBlockHeight, &dlg.Amount, &dlg.DelegationPeriod, &dlg.Created, &dlg.Started, &dlg.Finished, &dlg.Info, &dlg.State)
		if err != nil {
			return nil, err
		}
		delegations = append(delegations, dlg)
	}
	if len(delegations) == 0 {
		return nil, ErrNotFound
	}
	return delegations, nil
}

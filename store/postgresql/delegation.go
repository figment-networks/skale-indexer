package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	insertStatementD = `INSERT INTO delegations ("delegation_id", "holder", "validator_id", "amount", "delegation_period", "created", "started",  "finished", "info", "state" ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) `
	getByStatementD  = `SELECT d.id, d.created_at, d.updated_at, d.delegation_id, d.holder, d.validator_id, d.amount, d.delegation_period, d.created, d.started, d.finished, d.info, d.state FROM delegations d WHERE `
	byIdD            = `d.id =  $1 `
	byValidatorIdD   = `d.validator_id =  $1 `
	byDateRangeD     = `d.created between $1 and $2 `
	orderByCreatedD  = `ORDER BY d.created DESC`
)

// SaveDelegation saves delegation
func (d *Driver) SaveDelegation(ctx context.Context, dl structs.Delegation) error {
	_, err := d.db.Exec(insertStatementD, dl.DelegationID, dl.Holder, dl.ValidatorID, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info, dl.State)
	return err
}

// GetDelegations gets delegations by params
func (d *Driver) GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s%s", getByStatementD, byIdD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else if params.ValidatorId > 0 {
		q = fmt.Sprintf("%s%s%s", getByStatementD, byValidatorIdD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		q = fmt.Sprintf("%s%s%s", getByStatementD, byDateRangeD, orderByCreatedD)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	} else {
		// unexpected select query
		return delegations, handler.ErrMissingParameter
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.Delegation{}
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.DelegationID, &dlg.Holder, &dlg.ValidatorID, &dlg.Amount, &dlg.DelegationPeriod, &dlg.Created, &dlg.Started, &dlg.Finished, &dlg.Info, &dlg.State)
		if err != nil {
			return nil, err
		}
		delegations = append(delegations, dlg)
	}
	if len(delegations) == 0 {
		return nil, handler.ErrNotFound
	}
	return delegations, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	insertStatementForDelegation = `INSERT INTO delegations ("updated_at", "holder", "validator_id", "amount", "delegation_period", "created", "started",  "finished", "info" ) VALUES ( NOW(), $1, $2, $3, $4, $5, $6, $7, $8) `
	updateStatementForDelegation = `UPDATE delegations SET updated_at = NOW(), holder = $1, validator_id = $2, amount = $3, delegation_period = $4, created = $5, started = $6, finished = $7, info = $8  WHERE id = $9 `
	getByStatementForDelegation  = `SELECT d.id, d.created_at, d.updated_at, d.holder, d.validator_id, d.amount, d.delegation_period, d.created, d.started, d.finished, d.info FROM delegations d WHERE `
	byIdForDelegation            = `d.id =  $1 `
	byHolderForDelegation        = `d.holder =  $1 `
	byValidatorIdForDelegation   = `d.validator_id =  $1 `
	byDateRangeForDelegation     = `v.created between $1 and $2 `
	orderByCreated               = `ORDER BY created DESC`
)

func (d *Driver) saveOrUpdateDelegation(ctx context.Context, dl structs.Delegation) error {
	if dl.ID == "" {
		_, err := d.db.Exec(insertStatementForDelegation, dl.Holder, dl.ValidatorId, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info)
		return err
	}
	_, err := d.db.Exec(updateStatementForDelegation, dl.Holder, dl.ValidatorId, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info, dl.ID)
	return err
}

// SaveOrUpdateDelegations saves or updates delegations
func (d *Driver) SaveOrUpdateDelegations(ctx context.Context, dls []structs.Delegation) error {
	for _, dl := range dls {
		if err := d.saveOrUpdateDelegation(ctx, dl); err != nil {
			return err
		}
	}
	return nil
}

func (d *Driver) GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s%s", getByStatementForDelegation, byIdForDelegation, orderByCreated)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else if params.ValidatorId > 0 {
		q = fmt.Sprintf("%s%s%s", getByStatementForDelegation, byValidatorIdForDelegation, orderByCreated)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else if params.Holder != "" {
		q = fmt.Sprintf("%s%s%s", getByStatementForDelegation, byHolderForDelegation, orderByCreated)
		rows, err = d.db.QueryContext(ctx, q, params.Holder)
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		q = fmt.Sprintf("%s%s%s", getByStatementForDelegation, byDateRangeForDelegation, orderByCreated)
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
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.Holder, &dlg.ValidatorId, &dlg.Amount, &dlg.DelegationPeriod, &dlg.Created, &dlg.Started, &dlg.Finished, &dlg.Info)
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

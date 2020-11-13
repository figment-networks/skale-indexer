package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/figment-networks/skale-indexer/structs"
	"github.com/figment-networks/skale-indexer/types"
)

var (
	ErrNotFound = errors.New("record not found")
)

const (
	insertStatement = `INSERT INTO delegations ("holder", "validator_id", "amount", "delegation_period", "created", "started",  "finished", "info", "created_at", "updated_at" ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) `
	updateStatement = `UPDATE delegations SET holder = $1, validator_id = $2, amount = $3, delegation_period = $4, created = $5, started = $6, finished = $7, info = $8 , updated_at = NOW() WHERE id = $9 `
	getByStatement  = `SELECT * FROM delegations where `
	ById            = "id =  $1 "
	ByHolder        = "holder =  $1 "
	ByAmount        = "amount =  $1 "
)

// SaveOrUpdateDelegation saves or updates delegation
func (d *Driver) SaveOrUpdateDelegation(ctx context.Context, dl structs.Delegation) error {
	_, err := d.GetDelegationById(ctx, dl.ID)
	if err != nil {
		_, err = d.db.Exec(insertStatement, dl.Holder, dl.ValidatorId, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info)
	} else {
		_, err = d.db.Exec(insertStatement, dl.Holder, dl.ValidatorId, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info, dl.ID)

	}

	return nil
}

// SaveOrUpdateDelegations saves or updates delegations
func (d *Driver) SaveOrUpdateDelegations(ctx context.Context, dls []structs.Delegation) error {
	for _, dl := range dls {
		if err := d.SaveOrUpdateDelegation(ctx, dl); err != nil {
			return err
		}
	}
	return nil
}

// GetDelegationById gets delegation by id
func (d *Driver) GetDelegationById(ctx context.Context, id types.ID) (res structs.Delegation, err error) {
	dlg := structs.Delegation{}
	q := fmt.Sprintf("%s%s", getByStatement, ById)

	row := d.db.QueryRowContext(ctx, q, id)
	if row == nil {
		return res, ErrNotFound
	}

	err = row.Scan(&dlg.ID, &dlg.Holder, &dlg.ValidatorId, &dlg.Amount, &dlg.DelegationPeriod, &dlg.Created, &dlg.Started, &dlg.Finished, &dlg.Info, &dlg.CreatedAt, &dlg.UpdatedAt)
	if err == sql.ErrNoRows {
		return res, ErrNotFound
	}
	return res, err
}

// GetDelegationsByHolder gets delegations by holder
func (d *Driver) GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error) {
	q := fmt.Sprintf("%s%s", getByStatement, ByHolder)
	rows, err := d.db.QueryContext(ctx, q, holder)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("query error: %w", err)
	default:
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.Delegation{}
		err = rows.Scan(&dlg.ID, &dlg.Holder, &dlg.ValidatorId, &dlg.Amount, &dlg.DelegationPeriod, &dlg.Created, &dlg.Started, &dlg.Finished, &dlg.Info, &dlg.CreatedAt, &dlg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		delegations = append(delegations, dlg)
	}
	return delegations, nil
}

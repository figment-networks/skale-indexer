package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/figment-networks/skale-indexer/structs"
)

var (
	ErrNotFound = errors.New("record not found")
)

const (
	insertStatementForDelegation = `INSERT INTO delegations ("created_at", "updated_at", "holder", "validator_id", "amount", "delegation_period", "created", "started",  "finished", "info" ) VALUES ( NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8) `
	updateStatementForDelegation = `UPDATE delegations SET updated_at = NOW(), holder = $1, validator_id = $2, amount = $3, delegation_period = $4, created = $5, started = $6, finished = $7, info = $8  WHERE id = $9 `
	getByStatementForDelegation  = `SELECT d.id, d.created_at, d.updated_at, d.holder, d.validator_id, d.amount, d.delegation_period, d.created, d.started, d.finished, d.info FROM delegations d WHERE `
	byIdForDelegation            = `d.id =  $1 `
	byHolderForDelegation        = `d.holder =  $1 `
	byValidatorIdForDelegation   = `d.validator_id =  $1 `
	orderByCreated               = `ORDER BY created DESC`
)

func (d *Driver) saveOrUpdateDelegation(ctx context.Context, dl structs.Delegation) error {
	var err error
	if dl.ID == "" {
		_, err = d.db.Exec(insertStatementForDelegation, dl.Holder, dl.ValidatorId, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info)
	} else {
		_, err = d.db.Exec(updateStatementForDelegation, dl.Holder, dl.ValidatorId, dl.Amount, dl.DelegationPeriod, dl.Created, dl.Started, dl.Finished, dl.Info, dl.ID)
	}
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

// GetDelegationById gets delegation by id
func (d *Driver) GetDelegationById(ctx context.Context, id string) (res structs.Delegation, err error) {
	dlg := structs.Delegation{}
	q := fmt.Sprintf("%s%s", getByStatementForDelegation, byIdForDelegation)

	row := d.db.QueryRowContext(ctx, q, id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", row.Err().Error())
	}

	err = row.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.Holder, &dlg.ValidatorId, &dlg.Amount, &dlg.DelegationPeriod, &dlg.Created, &dlg.Started, &dlg.Finished, &dlg.Info)
	if err == sql.ErrNoRows || !(dlg.ID != "") {
		return res, ErrNotFound
	}
	return dlg, err
}

// GetDelegationsByHolder gets delegations by holder
func (d *Driver) GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error) {
	q := fmt.Sprintf("%s%s%s", getByStatementForDelegation, byHolderForDelegation, orderByCreated)
	rows, err := d.db.QueryContext(ctx, q, holder)
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
		return nil, ErrNotFound
	}
	return delegations, nil
}

// GetDelegationsByValidatorId gets delegations by validator id
func (d *Driver) GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error) {
	q := fmt.Sprintf("%s%s%s", getByStatementForDelegation, byValidatorIdForDelegation, orderByCreated)
	rows, err := d.db.QueryContext(ctx, q, validatorId)
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
		return nil, ErrNotFound
	}
	return delegations, nil
}

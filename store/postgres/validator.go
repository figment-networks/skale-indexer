package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
	"github.com/lib/pq"
)

const (
	insertStatementForValidator = `INSERT INTO validators ("updated_at", "name", "address", "description", "fee_rate", "active", "active_nodes", "staked", "pending", "rewards", "data") VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10) `
	updateStatementForValidator = `UPDATE validators SET updated_at = NOW(), name = $1, address = $2, description = $3, fee_rate = $4, active = $5, active_nodes = $6, staked = $7, pending = $8, rewards = $9, data = $10  WHERE id = $11 `
	getByStatementForValidator  = `SELECT v.id, v.created_at, v.updated_at, v.name, v.address, v.description, v.fee_rate, v.active, v.active_nodes, v.staked, v.pending, v.rewards, v.data FROM validators v WHERE `
	byIdForValidator            = `v.id =  $1 `
	byDateRange                 = `v.created_at between $1 and $2`
	byAddressForValidator       = `v.address =  $1 `
)

func (d *Driver) saveOrUpdateValidator(ctx context.Context, v structs.Validator) error {
	if v.ID == "" {
		_, err := d.db.Exec(insertStatementForValidator, v.Name, pq.Array(v.Address), v.Description, v.FeeRate, v.Active, v.ActiveNodes, v.Staked, v.Pending, v.Rewards, &v.OptionalInfo)
		return err
	}
	_, err := d.db.Exec(updateStatementForValidator, v.Name, pq.Array(v.Address), v.Description, v.FeeRate, v.Active, v.ActiveNodes, v.Staked, v.Pending, v.Rewards, &v.OptionalInfo, v.ID)
	return err
}

// SaveOrUpdateValidators saves or updates validators
func (d *Driver) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	for _, v := range validators {
		if err := d.saveOrUpdateValidator(ctx, v); err != nil {
			return err
		}
	}
	return nil
}

// GetValidatorById gets validator by id
func (d *Driver) GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error) {
	vld := structs.Validator{}
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byIdForValidator)

	row := d.db.QueryRowContext(ctx, q, id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", row.Err().Error())
	}

	err = row.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.Name, pq.Array(&vld.Address), &vld.Description, &vld.FeeRate, &vld.Active, &vld.ActiveNodes, &vld.Staked, &vld.Pending, &vld.Rewards, &vld.OptionalInfo)
	if err == sql.ErrNoRows || !(vld.ID != "") {
		return res, handler.ErrNotFound
	}
	return vld, err
}

// GetValidatorsByAddress gets validators by address
func (d *Driver) GetValidatorsByAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error) {
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byAddressForValidator)
	rows, err := d.db.QueryContext(ctx, q, validatorAddress)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.Name, pq.Array(&vld.Address), &vld.Description, &vld.FeeRate, &vld.Active, &vld.ActiveNodes, &vld.Staked, &vld.Pending, &vld.Rewards, &vld.OptionalInfo)
		if err != nil {
			return nil, err
		}
		validators = append(validators, vld)
	}
	if len(validators) == 0 {
		return nil, handler.ErrNotFound
	}
	return validators, nil
}

// GetValidators gets validators by params
func (d *Driver) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForValidator, byIdForValidator)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else {
		q = fmt.Sprintf("%s%s", getByStatementForValidator, byDateRange)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.Name, pq.Array(&vld.Address), &vld.Description, &vld.FeeRate, &vld.Active, &vld.ActiveNodes, &vld.Staked, &vld.Pending, &vld.Rewards, &vld.OptionalInfo)
		if err != nil {
			return nil, err
		}
		validators = append(validators, vld)
	}
	if len(validators) == 0 {
		return nil, handler.ErrNotFound
	}
	return validators, nil
}

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	insertStatementForValidator = `INSERT INTO validators ("updated_at", "validator_id", "name", "validator_address", "requested_address", "description", "fee_rate","registration_time", "minimum_delegation_amount", "accept_new_requests", "authorized", "active", "active_nodes", "linked_nodes", "staked", "pending", "rewards") VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) `
	getByStatementForValidator  = `SELECT v.id, v.created_at, v.updated_at, v.validator_id, v.name, v.validator_address, v.requested_address, v.description, v.fee_rate, v.registration_time, v.minimum_delegation_amount, v.accept_new_requests, v.authorized, v.active, v.active_nodes, v.linked_nodes, v.staked, v.pending, v.rewards FROM validators v WHERE `
	byIdForValidator            = `v.id =  $1 `
	byDateRangeForValidator     = `v.created_at between $1 and $2 `
	byAddressForValidator       = `v.validator_address =  $1 `
)

// SaveValidator saves validator
func (d *Driver) SaveValidator(ctx context.Context, v structs.Validator) error {
	_, err := d.db.Exec(insertStatementForValidator, v.ValidatorID, v.Name, v.ValidatorAddress, v.RequestedAddress, v.Description, v.FeeRate, v.RegistrationTime, v.MinimumDelegationAmount, v.AcceptNewRequests, v.Authorized, v.Active, v.ActiveNodes, v.LinkedNodes, v.Staked, v.Pending, v.Rewards)
	return err
}

// GetValidators gets validators by params
func (d *Driver) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForValidator, byIdForValidator)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	//} else if params.ValidatorAddress > 0 {
	//	q = fmt.Sprintf("%s%s", getByStatementForValidator, byAddressForValidator)
	//	rows, err = d.db.QueryContext(ctx, q, params.ValidatorAddress)
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		q = fmt.Sprintf("%s%s", getByStatementForValidator, byDateRangeForValidator)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	} else {
		// unexpected select query
		return validators, handler.ErrMissingParameter
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.ValidatorID, &vld.Name, &vld.ValidatorAddress, &vld.RequestedAddress, &vld.Description, &vld.FeeRate, &vld.RegistrationTime, &vld.MinimumDelegationAmount, &vld.AcceptNewRequests, &vld.Authorized, &vld.Active, &vld.ActiveNodes, &vld.LinkedNodes, &vld.Staked, &vld.Pending, &vld.Rewards)
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

package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// TODO: run explain analyze to check full scan and add required indexes
const (
	insertStatementV         = `INSERT INTO validators ("validator_id", "name", "validator_address", "requested_address", "description", "fee_rate","registration_time", "minimum_delegation_amount", "accept_new_requests", "authorized", "active", "active_nodes", "linked_nodes", "staked", "pending", "rewards", "block_height") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) `
	getByStatementV          = `SELECT v.id, v.created_at, v.updated_at, v.validator_id, v.name, v.validator_address, v.requested_address, v.description, v.fee_rate, v.registration_time, v.minimum_delegation_amount, v.accept_new_requests, v.authorized, v.active, v.active_nodes, v.linked_nodes, v.staked, v.pending, v.rewards, v.block_height FROM validators v WHERE `
	byIdV                    = `v.id =  $1 `
	byDateRangeV             = `v.created_at between $1 and $2 `
	byValidatorIdV           = `v.validator_id =  $1 `
	byRecentEthBlockHeightV  = `AND v.block_height =  (SELECT v2.block_height FROM delegations v2 WHERE v2.validator_id = $2 ORDER BY v2.block_height DESC LIMIT 1) `
	orderByRegistrationTimeV = `ORDER BY v.registration_time DESC `
)

// SaveValidator saves validator
func (d *Driver) SaveValidator(ctx context.Context, v structs.Validator) error {
	_, err := d.db.Exec(insertStatementV,
		v.ValidatorID.String(),
		v.Name,
		v.ValidatorAddress.Hash().Big().String(),
		v.RequestedAddress.Hash().Big().String(),
		v.Description,
		v.FeeRate.String(),
		v.RegistrationTime,
		v.MinimumDelegationAmount.String(),
		v.AcceptNewRequests,
		v.Authorized, v.Active, v.ActiveNodes, v.LinkedNodes, v.Staked, v.Pending, v.Rewards, v.BlockHeight)
	return err
}

// GetValidators gets validators by params
func (d *Driver) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	var q string
	var rows *sql.Rows
	if params.ValidatorId != 0 && !params.Recent {
		q = fmt.Sprintf("%s%s%s", getByStatementV, byValidatorIdV, orderByRegistrationTimeV)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else if params.ValidatorId != 0 && params.Recent {
		q = fmt.Sprintf("%s%s%s%s", getByStatementV, byValidatorIdV, byRecentEthBlockHeightV, orderByRegistrationTimeV)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId, params.ValidatorId)
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		q = fmt.Sprintf("%s%s%s", getByStatementV, byDateRangeV, orderByRegistrationTimeV)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	} else {
		q = fmt.Sprintf("%s%s", getByStatementV, byIdV)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.ValidatorID, &vld.Name, &vld.ValidatorAddress, &vld.RequestedAddress, &vld.Description, &vld.FeeRate, &vld.RegistrationTime, &vld.MinimumDelegationAmount, &vld.AcceptNewRequests, &vld.Authorized, &vld.Active, &vld.ActiveNodes, &vld.LinkedNodes, &vld.Staked, &vld.Pending, &vld.Rewards, &vld.ETHBlockHeight)
		if err != nil {
			return nil, err
		}
		validators = append(validators, vld)
	}
	return validators, nil
}

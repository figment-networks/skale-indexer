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
	insertStatementV         = `INSERT INTO validators ("validator_id", "name", "validator_address", "requested_address", "description", "fee_rate","registration_time", "minimum_delegation_amount", "accept_new_requests", "authorized", "active", "active_nodes", "linked_nodes", "staked", "pending", "rewards") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) `
	byValidatorIdV           = `AND validator_id =  $3 `
	orderByRegistrationTimeV = `ORDER BY registration_time DESC `
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
		v.Authorized, v.Active, v.ActiveNodes, v.LinkedNodes, v.Staked, v.Pending, v.Rewards)
	return err
}

// GetValidators gets validators by params
func (d *Driver) GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error) {
	// TODO: add minimum_delegation_amount
	q:= `SELECT id, created_at, validator_id, name, validator_address, requested_address, description, 
 				fee_rate, registration_time, accept_new_requests, authorized,
 				active, active_nodes, linked_nodes, staked, pending, rewards 
 		FROM validators WHERE created_at between $1 AND $2 `
	var rows *sql.Rows

	if params.ValidatorId != ""{
		q = fmt.Sprintf("%s%s", q, byValidatorIdV)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, params.ValidatorId)
	} else {
		q = fmt.Sprintf("%s%s", q, orderByRegistrationTimeV)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		var vldId uint64
		var validatorAddress []byte
		var requestedAddress []byte
		var feeRate uint64
		//var mnmDlgAmount uint64

		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vldId, &vld.Name, &validatorAddress, &requestedAddress, &vld.Description, &feeRate, &vld.RegistrationTime,
			//&mnmDlgAmount,
			&vld.AcceptNewRequests, &vld.Authorized, &vld.Active, &vld.ActiveNodes, &vld.LinkedNodes, &vld.Staked, &vld.Pending, &vld.Rewards)
		if err != nil {
			return nil, err
		}

		vld.ValidatorID = new(big.Int).SetUint64(vldId)
		vldAddress := new(big.Int)
		vldAddress.SetString(string(validatorAddress), 10)
		vld.ValidatorAddress.SetBytes(vldAddress.Bytes())
		rqtAddress := new(big.Int)
		rqtAddress.SetString(string(requestedAddress), 10)
		vld.RequestedAddress.SetBytes(rqtAddress.Bytes())
		vld.FeeRate = new(big.Int).SetUint64(feeRate)
		//vld.MinimumDelegationAmount = new(big.Int).SetUint64(mnmDlgAmount)

		validators = append(validators, vld)
	}
	return validators, nil
}

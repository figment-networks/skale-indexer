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
	byValidatorIdV       = `AND validator_id =  $3 `
	byRecentBlockHeightV = `AND block_height =  ((SELECT v2.block_height FROM validators v2 WHERE v2.validator_id = $3 ORDER BY v2.block_height DESC LIMIT 1)) `
	orderByValidatorIdV  = `ORDER BY validator_id `
)

// SaveValidator saves validator
func (d *Driver) SaveValidator(ctx context.Context, v structs.Validator) error {
	_, err := d.db.Exec(`INSERT INTO validators (
			"validator_id", 
			"name", 
			"validator_address", 
			"requested_address", 
			"description", 
			"fee_rate", 
			"block_height", 
			"registration_time", 
			"minimum_delegation_amount", 
			"accept_new_requests", 
			"authorized", 
			"active_nodes", 
			"linked_nodes", 
			"staked", 
			"pending", 
			"rewards") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) 
		ON CONFLICT (validator_id, block_height)
		DO UPDATE SET
			name = EXCLUDED.name,
			validator_address = EXCLUDED.validator_address,
			requested_address = EXCLUDED.requested_address,
			description = EXCLUDED.description,
			fee_rate = EXCLUDED.fee_rate,
			registration_time = EXCLUDED.registration_time,
			minimum_delegation_amount = EXCLUDED.minimum_delegation_amount,
			accept_new_requests = EXCLUDED.accept_new_requests,
			authorized = EXCLUDED.authorized,
			active_nodes = EXCLUDED.active_nodes,
			linked_nodes = EXCLUDED.linked_nodes,
			staked = EXCLUDED.staked,
			pending = EXCLUDED.pending,
			rewards = EXCLUDED.rewards
		`,
		v.ValidatorID.String(),
		v.Name,
		v.ValidatorAddress.Hash().Big().String(),
		v.RequestedAddress.Hash().Big().String(),
		v.Description,
		v.FeeRate.String(),
		v.BlockHeight,
		v.RegistrationTime,
		v.MinimumDelegationAmount.String(),
		v.AcceptNewRequests,
		v.Authorized, v.ActiveNodes, v.LinkedNodes, v.Staked, v.Pending, v.Rewards)
	return err
}

// GetValidators gets validators by params
func (d *Driver) GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error) {
	q := `SELECT id, created_at, validator_id, name, validator_address, requested_address, description, fee_rate, block_height, registration_time, 
			accept_new_requests, authorized, active_nodes, linked_nodes, staked, pending, rewards 
 		FROM validators WHERE created_at between $1 AND $2 `
	var rows *sql.Rows

	if params.ValidatorId != "" && params.Recent {
		q = fmt.Sprintf("%s%s%s", q, byValidatorIdD, byRecentBlockHeightV)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, params.ValidatorId)
	} else if params.ValidatorId != "" && !params.Recent {
		q = fmt.Sprintf("%s%s", q, byValidatorIdV)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, params.ValidatorId)
	} else if !params.Recent {
		q = fmt.Sprintf("%s%s", q, orderByValidatorIdV)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	} else {
		q = `SELECT DISTINCT ON (validator_id)  id, created_at, validator_id, name, validator_address, requested_address, description, fee_rate, block_height, registration_time, 
			accept_new_requests, authorized, active_nodes, linked_nodes, staked, pending, rewards 
 		FROM validators ORDER BY validator_id, block_height DESC`
		rows, err = d.db.QueryContext(ctx, q)
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
		//var mnmDlgAmount uint256.In

		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vldId, &vld.Name, &validatorAddress, &requestedAddress, &vld.Description, &feeRate, &vld.BlockHeight, &vld.RegistrationTime,
			//&mnmDlgAmount,
			&vld.AcceptNewRequests, &vld.Authorized, &vld.ActiveNodes, &vld.LinkedNodes, &vld.Staked, &vld.Pending, &vld.Rewards)
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
		// TODO: bug for some values out of range mnmDlgAmount
		//vld.MinimumDelegationAmount = new(big.Int).SetUint64(mnmDlgAmount)

		validators = append(validators, vld)
	}
	return validators, nil
}

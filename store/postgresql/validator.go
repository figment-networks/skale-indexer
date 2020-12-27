package postgresql

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
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
			"registration_time", 
			"minimum_delegation_amount", 
			"accept_new_requests", 
			"authorized", 
			"active_nodes", 
			"linked_nodes", 
			"staked", 
			"pending", 
			"rewards",
			"block_height") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) 
		ON CONFLICT (validator_id)
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
			pending = EXCLUDED.pending,
			rewards = EXCLUDED.rewards,
			block_height = EXCLUDED.block_height
		`,
		v.ValidatorID.String(),
		v.Name,
		v.ValidatorAddress.Hash().Big().String(),
		v.RequestedAddress.Hash().Big().String(),
		v.Description,
		v.FeeRate.String(),
		v.RegistrationTime,
		v.MinimumDelegationAmount.String(),
		v.AcceptNewRequests,
		v.Authorized,
		v.ActiveNodes,
		v.LinkedNodes,
		v.Staked.String(),
		v.Pending.String(),
		v.Rewards.String(),
		v.BlockHeight)
	return err
}

// GetValidators gets validators by params
func (d *Driver) GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error) {
	q := `SELECT id, created_at, validator_id, name, validator_address, requested_address, description, fee_rate, registration_time, minimum_delegation_amount, accept_new_requests, authorized, active_nodes, linked_nodes, staked, pending, rewards, block_height 
			FROM validators  `

	var (
		args   []interface{}
		wherec []string
		i      = 1
	)

	if params.ValidatorId != "" {
		wherec = append(wherec, ` validator_id =  $`+strconv.Itoa(i))
		args = append(args, params.ValidatorId)
		i++
	}
	if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		wherec = append(wherec, ` registration_time BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}

	if len(args) > 0 {
		q += ` WHERE `
	}
	q += strings.Join(wherec, " AND ")
	q += ` ORDER BY validator_id ASC`

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		var vldId uint64
		var validatorAddress []byte
		var requestedAddress []byte
		var feeRate uint64
		var mnmDlgAmount string
		var staked string
		var pending string
		var rewards string
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vldId, &vld.Name, &validatorAddress, &requestedAddress, &vld.Description, &feeRate, &vld.RegistrationTime, &mnmDlgAmount, &vld.AcceptNewRequests, &vld.Authorized, &vld.ActiveNodes, &vld.LinkedNodes, &staked, &pending, &rewards, &vld.BlockHeight)
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
		amnt := new(big.Int)
		amnt.SetString(mnmDlgAmount, 10)
		vld.MinimumDelegationAmount = amnt
		stk := new(big.Int)
		stk.SetString(staked, 10)
		vld.Staked = stk
		pnd := new(big.Int)
		pnd.SetString(pending, 10)
		vld.Pending = pnd
		rwrd := new(big.Int)
		rwrd.SetString(rewards, 10)
		vld.Rewards = rwrd
		validators = append(validators, vld)
	}
	return validators, nil
}

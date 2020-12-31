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

	if v.Staked == nil {
		v.Staked = big.NewInt(0)
	}
	if v.Pending == nil {
		v.Pending = big.NewInt(0)
	}
	if v.Rewards == nil {
		v.Rewards = big.NewInt(0)
	}

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
	q := `SELECT
				id,
				created_at,
				validator_id,
				name,
				validator_address,
				requested_address,
				description,
				fee_rate,
				registration_time,
				minimum_delegation_amount,
				accept_new_requests,
				authorized,
				active_nodes,
				linked_nodes,
				staked,
				pending,
				rewards,
				block_height
			FROM validators`

	var (
		args   []interface{}
		whereC []string
		i      = 1
	)

	if params.ValidatorId != "" {
		whereC = append(whereC, ` validator_id = $`+strconv.Itoa(i))
		args = append(args, params.ValidatorId)
		i++
	}
	if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		whereC = append(whereC, ` registration_time BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}

	if len(args) > 0 {
		q += ` WHERE `
	}
	q += strings.Join(whereC, " AND ")
	q += ` ORDER BY validator_id ASC`

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		vldID        uint64
		feeRate      uint64
		mnmDlgAmount uint64
		staked       uint64
		pending      uint64
		rewards      uint64
	)
	for rows.Next() {
		vld := structs.Validator{}
		var validatorAddress []byte
		var requestedAddress []byte
		err = rows.Scan(&vld.ID,
			&vld.CreatedAt,
			&vldID,
			&vld.Name,
			&validatorAddress,
			&requestedAddress,
			&vld.Description,
			&feeRate,
			&vld.RegistrationTime,
			&mnmDlgAmount,
			&vld.AcceptNewRequests,
			&vld.Authorized,
			&vld.ActiveNodes,
			&vld.LinkedNodes,
			&staked,
			&pending,
			&rewards,
			&vld.BlockHeight)
		if err != nil {
			return nil, err
		}
		vld.ValidatorID = new(big.Int).SetUint64(vldID)
		vInt := new(big.Int)
		vInt.SetString(string(requestedAddress), 10)
		vld.RequestedAddress.SetBytes(vInt.Bytes())
		vInt.SetString(string(validatorAddress), 10)
		vld.ValidatorAddress.SetBytes(vInt.Bytes())

		vld.FeeRate = new(big.Int).SetUint64(feeRate)
		vld.MinimumDelegationAmount = new(big.Int).SetUint64(mnmDlgAmount)
		vld.Staked = new(big.Int).SetUint64(staked)
		vld.Pending = new(big.Int).SetUint64(pending)
		vld.Rewards = new(big.Int).SetUint64(rewards)
		validators = append(validators, vld)
	}
	return validators, nil
}

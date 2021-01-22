package postgresql

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

var zerobig = big.NewInt(0)

// SaveValidators saves validators
func (d *Driver) SaveValidators(ctx context.Context, validators []structs.Validator) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, v := range validators {
		if v.Staked == nil {
			v.Staked = zerobig
		}
		if v.FeeRate == nil {
			v.FeeRate = zerobig
		}
		if v.MinimumDelegationAmount == nil {
			v.MinimumDelegationAmount = zerobig
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO validators (
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
			"block_height")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
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
			v.BlockHeight)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	return tx.Commit()
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
				block_height
			FROM validators`

	var (
		args   []interface{}
		whereC []string
		i      = 1
	)

	if params.ValidatorID != "" {
		whereC = append(whereC, ` validator_id = $`+strconv.Itoa(i))
		args = append(args, params.ValidatorID)
		i++
	}
	if params.Active != "" {
		whereC = append(whereC, ` authorized =  $`+strconv.Itoa(i))
		args = append(args, params.Active)
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
	q += ` ORDER BY `
	if params.OrderBy != "" {
		q += params.OrderBy
		if params.OrderDirection != "" {
			q += ` ` + params.OrderDirection
		}
	} else {
		q += ` validator_id ASC `
	}

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		vldID        uint64
		feeRate      string
		mnmDlgAmount string
		staked       string
	)
	for rows.Next() {
		vld := structs.Validator{}
		var (
			validatorAddress []byte
			requestedAddress []byte
		)
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

		vld.FeeRate, _ = new(big.Int).SetString(feeRate, 10)
		vld.MinimumDelegationAmount, _ = new(big.Int).SetString(mnmDlgAmount, 10)
		vld.Staked, _ = new(big.Int).SetString(staked, 10)
		validators = append(validators, vld)
	}
	return validators, nil
}

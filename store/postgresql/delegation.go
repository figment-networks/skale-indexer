package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// SaveDelegations saves delegations
func (d *Driver) SaveDelegations(ctx context.Context, delegations []structs.Delegation) error {

	for _, dl := range delegations {
		_, err := d.db.Exec(`INSERT INTO delegations (
				"delegation_id",
				"holder",
				"validator_id",
				"block_height",
				"transaction_hash",
				"amount",
				"delegation_period",
				"created",
				"started",
				"finished",
				"info",
				"state")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (delegation_id, transaction_hash)
		DO UPDATE SET
			holder = EXCLUDED.holder,
			block_height = EXCLUDED.block_height,
			validator_id = EXCLUDED.validator_id,
			amount = EXCLUDED.amount,
			delegation_period = EXCLUDED.delegation_period,
			created = EXCLUDED.created,
			started = EXCLUDED.started,
			finished = EXCLUDED.finished,
			info = EXCLUDED.info,
			state = EXCLUDED.state
		`,
			dl.DelegationID.String(),
			dl.Holder.Hash().Big().String(),
			dl.ValidatorID.String(),
			dl.BlockHeight,
			dl.TransactionHash.Big().String(),
			dl.Amount.String(),
			dl.DelegationPeriod.String(),
			dl.Created,
			dl.Started.String(),
			dl.Finished.String(),
			dl.Info,
			dl.State)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetDelegationTimeline gets all delegation information over time
func (d *Driver) GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	q := `SELECT id, delegation_id, holder, validator_id, block_height, transaction_hash, amount, delegation_period, created, started, finished, info, state
			FROM delegations `

	var (
		args   []interface{}
		whereC []string
		i      = 1
	)

	if params.DelegationID != "" {
		whereC = append(whereC, ` delegation_id = $`+strconv.Itoa(i))
		args = append(args, params.DelegationID)
		i++
	}
	if params.ValidatorID != "" {
		whereC = append(whereC, ` validator_id = $`+strconv.Itoa(i))
		args = append(args, params.ValidatorID)
		i++
	}
	if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		whereC = append(whereC, ` created BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}

	if len(whereC) > 0 {
		q += " WHERE "
	}
	q += strings.Join(whereC, " AND ")
	q += `ORDER BY block_height DESC`

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dlg := structs.Delegation{}
		var (
			th        []byte
			dlgID     uint64
			holder    []byte
			vldID     uint64
			amount    string
			started   uint64
			finished  uint64
			dlgPeriod uint64
		)

		if err := rows.Scan(&dlg.ID, &dlgID, &holder, &vldID, &dlg.BlockHeight, &th, &amount, &dlgPeriod, &dlg.Created, &started, &finished, &dlg.Info, &dlg.State); err != nil {
			return nil, err
		}

		h := new(big.Int)
		h.SetString(string(holder), 10)
		dlg.Holder.SetBytes(h.Bytes())

		h.SetString(string(th), 10)
		dlg.TransactionHash.SetBytes(h.Bytes())

		h.SetString(amount, 10)
		dlg.Amount = h

		dlg.ValidatorID = new(big.Int).SetUint64(vldID)
		dlg.DelegationID = new(big.Int).SetUint64(dlgID)
		dlg.DelegationPeriod = new(big.Int).SetUint64(dlgPeriod)
		dlg.Started = new(big.Int).SetUint64(started)
		dlg.Finished = new(big.Int).SetUint64(finished)
		delegations = append(delegations, dlg)
	}
	return delegations, nil
}

func (d *Driver) GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	q := `SELECT
			DISTINCT ON (delegation_id)
				delegation_id, id, holder, validator_id, block_height, transaction_hash, amount, delegation_period, created, started, finished, info, state
			FROM delegations `

	var (
		args   []interface{}
		whereC []string
		i      = 1
	)

	if params.DelegationID != "" {
		whereC = append(whereC, ` delegation_id = $`+strconv.Itoa(i))
		args = append(args, params.DelegationID)
		i++
	}
	if params.ValidatorID != "" {
		whereC = append(whereC, ` validator_id = $`+strconv.Itoa(i))
		args = append(args, params.ValidatorID)
		i++
	}
	if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		whereC = append(whereC, ` created BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}
	if len(whereC) > 0 {
		q += " WHERE "
	}
	q += strings.Join(whereC, " AND ")
	q += `ORDER BY delegation_id DESC, block_height DESC`

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dlg := structs.Delegation{}
		var (
			th        []byte
			dlgId     uint64
			holder    []byte
			vldId     uint64
			amount    string
			started   uint64
			finished  uint64
			dlgPeriod uint64
		)

		if err := rows.Scan(&dlgId, &dlg.ID, &holder, &vldId, &dlg.BlockHeight, &th, &amount, &dlgPeriod, &dlg.Created, &started, &finished, &dlg.Info, &dlg.State); err != nil {
			return nil, err
		}

		h := new(big.Int)
		h.SetString(string(holder), 10)
		dlg.Holder.SetBytes(h.Bytes())

		h.SetString(string(th), 10)
		dlg.TransactionHash.SetBytes(h.Bytes())

		h.SetString(amount, 10)
		dlg.Amount = h

		dlg.ValidatorID = new(big.Int).SetUint64(vldId)
		dlg.DelegationID = new(big.Int).SetUint64(dlgId)
		dlg.DelegationPeriod = new(big.Int).SetUint64(dlgPeriod)
		dlg.Started = new(big.Int).SetUint64(started)
		dlg.Finished = new(big.Int).SetUint64(finished)
		delegations = append(delegations, dlg)
	}
	return delegations, nil
}

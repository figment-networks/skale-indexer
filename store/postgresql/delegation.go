package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// SaveDelegation saves delegation
func (d *Driver) SaveDelegation(ctx context.Context, dl structs.Delegation) error {

	_, err := d.db.Exec(`INSERT INTO delegations (
				"delegation_id","holder","validator_id", "block_height","transaction_hash","amount",
				"delegation_period","created","started","finished","info","state","until")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, (date_trunc('month', $14::timestamp ) + make_interval(MONTHS => 1+$13::INTEGER) - interval '1 day')::TIMESTAMP )
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
			state = EXCLUDED.state,
			until = EXCLUDED.until
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
		dl.State,
		dl.DelegationPeriod.Uint64(),
		dl.Created.Format("2006-01-02 15:04:05"))

	return err
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
	if params.Holder != "" {
		whereC = append(whereC, ` holder =  $`+strconv.Itoa(i))
		args = append(args, common.HexToAddress(params.Holder).Hash().Big().String())
		i++
	}

	if !params.TimeAt.IsZero() {
		whereC = append(whereC, `$`+strconv.Itoa(i)+` BETWEEN created AND until`)
		args = append(args, params.TimeAt)
		i += 1
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		whereC = append(whereC, ` created BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}
	if len(params.State) > 0 {
		whereC = append(whereC, "state @> $"+strconv.Itoa(i))
		args = append(args, pq.Array(params.State))
		i++
	}

	if len(whereC) > 0 {
		q += " WHERE "
	}
	q += strings.Join(whereC, " AND ")
	q += `ORDER BY block_height DESC`

	if params.Limit > 0 {
		q += " LIMIT " + strconv.FormatUint(uint64(params.Limit), 10)
		if params.Offset > 0 {
			q += " OFFSET " + strconv.FormatUint(uint64(params.Offset), 10)
		}
	}
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
	if params.Holder != "" {
		whereC = append(whereC, ` holder =  $`+strconv.Itoa(i))
		args = append(args, common.HexToAddress(params.Holder).Hash().Big().String())
		i++
	}

	if len(params.State) > 0 {
		whereC = append(whereC, "state @> $"+strconv.Itoa(i))
		args = append(args, pq.Array(params.State))
		i++
	}

	if !params.TimeAt.IsZero() {
		whereC = append(whereC, `$`+strconv.Itoa(i)+` BETWEEN created AND until`)
		args = append(args, params.TimeAt)
		i += 1
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		whereC = append(whereC, ` created BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}

	if len(whereC) > 0 {
		q += " WHERE "
	}
	q += strings.Join(whereC, " AND ")
	q += ` ORDER BY delegation_id DESC, block_height DESC`

	if params.Limit > 0 {
		q += " LIMIT " + strconv.FormatUint(uint64(params.Limit), 10)
		if params.Offset > 0 {
			q += " OFFSET " + strconv.FormatUint(uint64(params.Offset), 10)
		}
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

func (d *Driver) GetTypesSummaryDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.DelegationSummary, err error) {
	q := `SELECT DISTINCT ON (delegation_id) delegation_id, amount, state
				FROM delegations `

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

	if !params.TimeAt.IsZero() {
		whereC = append(whereC, `$`+strconv.Itoa(i)+` BETWEEN created AND until`)
		args = append(args, params.TimeAt)
		i += 1
	} else if !params.TimeFrom.IsZero() && !params.TimeTo.IsZero() {
		whereC = append(whereC, ` created BETWEEN $`+strconv.Itoa(i)+` AND $`+strconv.Itoa(i+1))
		args = append(args, params.TimeFrom)
		args = append(args, params.TimeTo)
		i += 2
	}
	if len(whereC) > 0 {
		q += " WHERE "
	}
	q += strings.Join(whereC, " AND ")
	q += ` ORDER BY delegation_id DESC, block_height DESC`

	var rows *sql.Rows

	rows, err = d.db.QueryContext(ctx, `SELECT d.state, SUM(d.amount) as amount, COUNT(d.delegation_id) as count FROM (`+q+`) AS d GROUP BY d.state;`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		amount string
		count  string
	)

	for rows.Next() {
		dlg := structs.DelegationSummary{}
		if err := rows.Scan(&dlg.State, &amount, &count); err != nil {
			return nil, err
		}

		h := new(big.Int)

		h.SetString(count, 10)
		dlg.Count = h

		h.SetString(amount, 10)
		dlg.Amount = h

		delegations = append(delegations, dlg)
	}
	return delegations, nil
}

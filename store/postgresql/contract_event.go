package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/lib/pq"
)

// SaveEvent saves contract events
func (d *Driver) SaveContractEvent(ctx context.Context, ce structs.ContractEvent) error {

	a := pq.Array([]string{ce.BoundID.String()})
	params, err := json.Marshal(ce.Params)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(
		`INSERT INTO contract_events(
			"contract_name",
			"event_name",
			"contract_address",
			"block_height",
			"time",
			"transaction_hash",
			"params",
			"removed",
			"bound_type",
			"bound_id",
			"bound_address")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		ce.ContractName,
		ce.EventName,
		ce.ContractAddress.Hash().Big().String(),
		ce.BlockHeight,
		ce.Time,
		ce.TransactionHash.Big().String(),
		params,
		ce.Removed,
		ce.BoundType,
		a,
		pq.Array([]string{ce.BoundAddress.Hash().Big().String()}),
	)
	return err
}

// GetContractEvents gets contract events
func (d *Driver) GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error) {

	q := `SELECT id, contract_name, event_name, contract_address, block_height, time, transaction_hash, params, removed
		FROM contract_events WHERE time BETWEEN $1 AND $2 `

	if params.Id > 0 {
		q += ` AND bound_id = $3 AND bound_type = $4`
	}

	q += ` ORDER BY time DESC`
	/*
		if params.BoundType == "validator" {
			q = fmt.Sprintf("%s%s%s%s", getByStatementCE, byTimeRangeCE, byBoundIdCE, orderByEventTimeCE)
		} else if params.BoundType == "delegation" {
			q = fmt.Sprintf("%s%s%s%s", getByStatementCE, byTimeRangeCE, byBoundIdCE, orderByEventTimeCE)
			rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, pq.Array(params.BoundId))
		}
	*/

	var rows *sql.Rows
	if params.Id > 0 {
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, pq.Array(params.Id), params.Type)
	} else {
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		e := structs.ContractEvent{}
		var ca []byte
		var th []byte
		var params []byte
		if err = rows.Scan(&e.ID, &e.ContractName, &e.EventName, &ca, &e.BlockHeight, &e.Time, &th, &params, &e.Removed); err != nil {
			return nil, err
		}
		p := new(big.Int)
		p.SetString(string(ca), 10)
		e.ContractAddress.SetBytes(p.Bytes())

		p.SetString(string(th), 10)
		e.TransactionHash.SetBytes(p.Bytes())

		a := make(map[string]interface{})
		if err := json.Unmarshal(params, &a); err != nil {
			return nil, fmt.Errorf("unmarshal error error: %w", err)
		}

		e.Params = a
		contractEvents = append(contractEvents, e)
	}

	return contractEvents, nil
}

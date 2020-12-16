package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/lib/pq"
)

// SaveEvent saves contract events
func (d *Driver) SaveContractEvent(ctx context.Context, ce structs.ContractEvent) error {
	params, err := json.Marshal(ce.Params)
	if err != nil {
		return err
	}

	var (
		bIDs   []string
		bAddrs []string
	)

	for _, bid := range ce.BoundID {
		bIDs = append(bIDs, bid.String())
	}
	for _, baddr := range ce.BoundAddress {
		bAddrs = append(bAddrs, baddr.Hash().Big().String())
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (contract_address, event_name, block_height, transaction_hash, removed)
		DO UPDATE SET
			contract_name = EXCLUDED.contract_name,
			time = EXCLUDED.time,
			params = EXCLUDED.params,
			bound_type = EXCLUDED.bound_type,
			bound_id = EXCLUDED.bound_id,
			bound_address = EXCLUDED.bound_address
		`,
		ce.ContractName,
		ce.EventName,
		ce.ContractAddress.Hash().Big().String(),
		ce.BlockHeight,
		ce.Time,
		ce.TransactionHash.Big().String(),
		params,
		ce.Removed,
		ce.BoundType,
		pq.Array(bIDs),
		pq.Array(bAddrs),
	)
	return err
}

// GetContractEvents gets contract events
func (d *Driver) GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error) {

	q := `SELECT id, contract_name, event_name, contract_address, block_height, time, transaction_hash, params, removed
		FROM contract_events WHERE time BETWEEN $1 AND $2 `

	if params.Type != "" {
		switch params.Type {
		case "validator":
			q += ` AND ( (bound_type = 'validator' AND bound_id[1] = $3 ) OR
				   ( bound_type = 'delegation' AND bound_id[2] = $3) ) `
		case "delegation":
			q += ` AND (bound_id[1] = $3 AND bound_type = 'delegation') `
		case "node":
			q += ` AND (bound_id[1] = $3 AND bound_type = 'node') `
		default:
			return nil, errors.New("unknown type")
		}
	}

	q += ` ORDER BY time DESC`

	var rows *sql.Rows
	if params.Id > 0 {
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, params.Id)
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

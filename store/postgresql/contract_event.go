package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/lib/pq"
)

const (
	insertStatementCE        = `INSERT INTO contract_events ("contract_name", "event_name", "contract_address", "block_height", "time", "transaction_hash", "params", "removed", "bound_type", "bound_id", "bound_address") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) `
	getByStatementCE         = `SELECT e.id, e.created_at, e.contract_name, e.event_name, e.contract_address, e.block_height, e.time, e.transaction_hash, e.params, e.removed, e.bound_type, e.bound_id, e.bound_address FROM contract_events e `
	byTimeRangeCE            = `WHERE e.time BETWEEN $1 AND $2 `
	byBoundIdCE              = `AND e.bound_id = $3 `
	orderByEventTimeCE       = `ORDER BY e.time DESC`
)

// SaveEvent saves contract events
func (d *Driver) SaveContractEvent(ctx context.Context, ce structs.ContractEvent) error {
	_, err := d.db.Exec(insertStatementCE, ce.ContractName, ce.EventName, pq.Array(ce.ContractAddress), ce.BlockHeight, ce.Time, ce.TransactionHash, ce.Params, ce.Removed, ce.BoundType, pq.Array(ce.BoundId), pq.Array(ce.BoundAddress))
	return err
}

// GetContractEvents gets contract events
func (d *Driver) GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error) {
	var q string
	var rows *sql.Rows

	if params.BoundType == "validator" {
		q = fmt.Sprintf("%s%s%s%s", getByStatementCE, byTimeRangeCE, byBoundIdCE, orderByEventTimeCE)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, pq.Array(params.BoundId))
	} else if params.BoundType == "delegation" {
		q = fmt.Sprintf("%s%s%s%s", getByStatementCE, byTimeRangeCE, byBoundIdCE, orderByEventTimeCE)
		rows, err = d.db.QueryContext(ctx, q, params.TimeFrom, params.TimeTo, pq.Array(params.BoundId))
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		e := structs.ContractEvent{}
		err = rows.Scan(&e.ID, &e.CreatedAt, &e.ContractName, &e.EventName, pq.Array(&e.ContractAddress), &e.BlockHeight, &e.Time, &e.TransactionHash, &e.Params, &e.Removed, &e.BoundType, &e.BoundId, pq.Array(&e.BoundAddress))
		if err != nil {
			return nil, err
		}
		contractEvents = append(contractEvents, e)
	}
	if len(contractEvents) == 0 {
		return nil, ErrNotFound
	}
	return contractEvents, nil
}

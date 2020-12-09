package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/client/structures"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	insertStatementForEvent = `INSERT INTO events ("block_height", "smart_contract_address", "transaction_index", "event_type", "event_name", "event_time", "event_info") VALUES ($1, $2, $3, $4, $5, $6) `
	updateStatementForEvent = `UPDATE events SET updated_at = NOW(), block_height = $1, smart_contract_address = $2, transaction_index = $3, event_type = $4, event_name = $5, event_time = $6, event_info = $7 WHERE id = $8 `
	getByStatementForEvent  = `SELECT e.id, e.created_at, e.updated_at, e.block_height, e.smart_contract_address, e.transaction_index, e.event_type, e.event_name, e.event_time, e.event_info FROM events e `
	byIdForEvent            = `WHERE e.id =  $1 `
	orderByEventTime        = `ORDER BY e.event_time DESC`
)

func (d *Driver) saveOrUpdateEvent(ctx context.Context, dl structs.Event) error {
	if dl.ID == "" {
		_, err := d.db.Exec(insertStatementForEvent, dl.BlockHeight, dl.SmartContractAddress, dl.TransactionIndex, dl.EventType, dl.EventName, dl.EventTime, dl.EventInfo)
		return err
	}
	_, err := d.db.Exec(updateStatementForEvent, dl.BlockHeight, dl.SmartContractAddress, dl.TransactionIndex, dl.EventType, dl.EventName, dl.EventTime, dl.EventInfo, dl.ID)
	return err
}

// SaveOrUpdateEvents saves or updates events
func (d *Driver) SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error {
	for _, dl := range events {
		if err := d.saveOrUpdateEvent(ctx, dl); err != nil {
			return err
		}
	}
	return nil
}

// GetEvents gets  events
func (d *Driver) GetEvents(ctx context.Context, params structs.QueryParams) (events []structs.Event, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForEvent, byIdForEvent)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else {
		q = fmt.Sprintf("%s%s", getByStatementForEvent, orderByEventTime)
		rows, err = d.db.QueryContext(ctx, q)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		e := structs.Event{}
		err = rows.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.BlockHeight, &e.SmartContractAddress, &e.TransactionIndex, &e.EventType, &e.EventName, &e.EventTime, &e.EventInfo)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if len(events) == 0 {
		return nil, handler.ErrNotFound
	}
	return events, nil
}

func (d *Driver) StoreEvent(ctx context.Context, boundAddress common.Address, boundType string, ev structures.ContractEvent) error {

	_, err := d.db.Exec(`INSERT INTO events("block_height", "event_time", "smart_contract_address", "event_type", "event_name", "transaction_hash", "bound_address", "bound_type", "event_info") VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		ev.Height,
		ev.Time,
		ev.Address,
		ev.Type,
		ev.ContractName,
		ev.TxHash,
		boundAddress,
		boundType,
		ev.Params)

	return err
}

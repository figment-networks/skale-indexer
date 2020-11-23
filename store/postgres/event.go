package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	insertStatementForEvent = `INSERT INTO events ("created_at", "updated_at", "block_height", "smart_contract_address", "transaction_index", "event_type", "event_name", "event_time") VALUES ( NOW(), NOW(), $1, $2, $3, $4, $5, $6) `
	updateStatementForEvent = `UPDATE events SET updated_at = NOW(), block_height = $1, smart_contract_address = $2, transaction_index = $3, event_type = $4, event_name = $5, event_time = $6 WHERE id = $7 `
	getByStatementForEvent  = `SELECT e.id, e.created_at, e.updated_at, e.block_height, e.smart_contract_address, e.transaction_index, e.event_type, e.event_name, e.event_time FROM events e `
	byIdForEvent            = `WHERE e.id =  $1 `
	orderByEventTime        = `ORDER BY e.event_time DESC`
)

func (d *Driver) saveOrUpdateEvent(ctx context.Context, dl structs.Event) error {
	if dl.ID == "" {
		_, err := d.db.Exec(insertStatementForEvent, dl.BlockHeight, dl.SmartContractAddress, dl.TransactionIndex, dl.EventType, dl.EventName, dl.EventTime)
		return err
	}
	_, err := d.db.Exec(updateStatementForEvent, dl.BlockHeight, dl.SmartContractAddress, dl.TransactionIndex, dl.EventType, dl.EventName, dl.EventTime, dl.ID)
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

// GetEventById gets event by id
func (d *Driver) GetEventById(ctx context.Context, id string) (res structs.Event, err error) {
	dlg := structs.Event{}
	q := fmt.Sprintf("%s%s", getByStatementForEvent, byIdForEvent)

	row := d.db.QueryRowContext(ctx, q, id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", row.Err().Error())
	}

	err = row.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.BlockHeight, &dlg.SmartContractAddress, &dlg.TransactionIndex, &dlg.EventType, &dlg.EventName, &dlg.EventTime)
	if err == sql.ErrNoRows || !(dlg.ID != "") {
		return res, ErrNotFound
	}
	return dlg, err
}

// GetAllEvents gets all events
func (d *Driver) GetAllEvents(ctx context.Context) (events []structs.Event, err error) {
	q := fmt.Sprintf("%s%s", getByStatementForEvent, orderByEventTime)
	rows, err := d.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.Event{}
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.BlockHeight, &dlg.SmartContractAddress, &dlg.TransactionIndex, &dlg.EventType, &dlg.EventName, &dlg.EventTime)
		if err != nil {
			return nil, err
		}
		events = append(events, dlg)
	}
	if len(events) == 0 {
		return nil, ErrNotFound
	}
	return events, nil
}

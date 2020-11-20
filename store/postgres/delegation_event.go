package postgres

import (
	"../../structs"
	"context"
	"database/sql"
	"fmt"
)

const (
	insertStatementForDelegationEvent = `INSERT INTO delegation_events ("created_at", "updated_at", "delegation_id", "event_name", "event_date") VALUES ( NOW(), NOW(), $1, $2, $3) `
	updateStatementForDelegationEvent = `UPDATE delegation_events SET updated_at = NOW(), delegation_id = $1, event_name = $2, event_date = $3 WHERE id = $4 `
	getByStatementForDelegationEvent  = `SELECT d.id, d.created_at, d.updated_at, d.delegation_id, d.event_name, d.event_date FROM delegation_events d `
	byIdForDelegationEvent            = `WHERE d.id =  $1 `
	byDelegationIdForDelegationEvent  = `WHERE d.delegation_id =  $1 `
	orderByEventTime                  = `ORDER BY event_time DESC`
)

func (d *Driver) saveOrUpdateDelegationEvent(ctx context.Context, dl structs.DelegationEvent) error {
	var err error
	if dl.ID == nil {
		_, err = d.db.Exec(insertStatementForDelegationEvent, dl.DelegationId, dl.EventName, dl.EventTime)
	} else {
		_, err = d.db.Exec(updateStatementForDelegationEvent, dl.DelegationId, dl.EventName, dl.EventTime, dl.ID)
	}
	return err
}

// SaveOrUpdateDelegationEvents saves or updates delegation events
func (d *Driver) SaveOrUpdateDelegationEvents(ctx context.Context, delegationEvents []structs.DelegationEvent) error {
	for _, dl := range delegationEvents {
		if err := d.saveOrUpdateDelegationEvent(ctx, dl); err != nil {
			return err
		}
	}
	return nil
}

// GetDelegationEventById gets delegation event by id
func (d *Driver) GetDelegationEventById(ctx context.Context, id *string) (res structs.DelegationEvent, err error) {
	dlg := structs.DelegationEvent{}
	q := fmt.Sprintf("%s%s", getByStatementForDelegationEvent, byIdForDelegationEvent)

	row := d.db.QueryRowContext(ctx, q, *id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", row.Err().Error())
	}

	err = row.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.DelegationId, &dlg.EventName, &dlg.EventTime)
	if err == sql.ErrNoRows || !(*dlg.ID != "") {
		return res, ErrNotFound
	}
	return dlg, err
}

// GetDelegationEventsByDelegationId gets delegation events by delegation id
func (d *Driver) GetDelegationEventsByDelegationId(ctx context.Context, delegationId *string) (delegationEvents []structs.DelegationEvent, err error) {
	q := fmt.Sprintf("%s%s%s", getByStatementForDelegationEvent, byDelegationIdForDelegationEvent, orderByEventTime)
	rows, err := d.db.QueryContext(ctx, q, *delegationId)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.DelegationEvent{}
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.DelegationId, &dlg.EventName, &dlg.EventTime)
		if err != nil {
			return nil, err
		}
		delegationEvents = append(delegationEvents, dlg)
	}
	if len(delegationEvents) == 0 {
		return nil, ErrNotFound
	}
	return delegationEvents, nil
}

// GetAllDelegationEvents gets all delegation events
func (d *Driver) GetAllDelegationEvents(ctx context.Context) (delegationEvents []structs.DelegationEvent, err error) {
	q := fmt.Sprintf("%s%s", getByStatementForDelegationEvent, orderByEventTime)
	rows, err := d.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		dlg := structs.DelegationEvent{}
		err = rows.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.DelegationId, &dlg.EventName, &dlg.EventTime)
		if err != nil {
			return nil, err
		}
		delegationEvents = append(delegationEvents, dlg)
	}
	if len(delegationEvents) == 0 {
		return nil, ErrNotFound
	}
	return delegationEvents, nil
}

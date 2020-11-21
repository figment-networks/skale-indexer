package postgres

import (
	"../../structs"
	"context"
	"database/sql"
	"fmt"
)

const (
	insertStatementForValidatorEvent = `INSERT INTO validator_events ("created_at", "updated_at", "validator_id", "event_name", "event_time") VALUES ( NOW(), NOW(), $1, $2, $3) `
	updateStatementForValidatorEvent = `UPDATE validator_events SET updated_at = NOW(), validator_id = $1, event_name = $2, event_time = $3 WHERE id = $4 `
	getByStatementForValidatorEvent  = `SELECT v.id, v.created_at, v.updated_at, v.validator_id, v.event_name, v.event_time FROM validator_events v `
	byIdForValidatorEvent            = `WHERE v.id =  $1 `
	byValidatorIdForValidatorEvent   = `WHERE v.validator_id =  $1 `
)

func (d *Driver) saveOrUpdateValidatorEvent(ctx context.Context, ve structs.ValidatorEvent) error {
	var err error
	if ve.ID == "" {
		_, err = d.db.Exec(insertStatementForValidatorEvent, ve.ValidatorId, ve.EventName, ve.EventTime)
	} else {
		_, err = d.db.Exec(updateStatementForValidatorEvent, ve.ValidatorId, ve.EventName, ve.EventTime, ve.ID)
	}
	return err
}

// SaveOrUpdateValidatorEvents saves or updates validator events
func (d *Driver) SaveOrUpdateValidatorEvents(ctx context.Context, validatorEvents []structs.ValidatorEvent) error {
	for _, ve := range validatorEvents {
		if err := d.saveOrUpdateValidatorEvent(ctx, ve); err != nil {
			return err
		}
	}
	return nil
}

// GetValidatorEventById gets validator event by id
func (d *Driver) GetValidatorEventById(ctx context.Context, id string) (res structs.ValidatorEvent, err error) {
	ve := structs.ValidatorEvent{}
	q := fmt.Sprintf("%s%s", getByStatementForValidatorEvent, byIdForValidatorEvent)

	row := d.db.QueryRowContext(ctx, q, id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", row.Err().Error())
	}

	err = row.Scan(&ve.ID, &ve.CreatedAt, &ve.UpdatedAt, &ve.ValidatorId, &ve.EventName, &ve.EventTime)
	if err == sql.ErrNoRows || !(ve.ID != "") {
		return res, ErrNotFound
	}
	return ve, err
}

// GetValidatorEventsByValidatorId gets validator events by validator id
func (d *Driver) GetValidatorEventsByValidatorId(ctx context.Context, validatorId string) (validatorEvents []structs.ValidatorEvent, err error) {
	q := fmt.Sprintf("%s%s%s", getByStatementForValidatorEvent, byValidatorIdForValidatorEvent, orderByEventTime)
	rows, err := d.db.QueryContext(ctx, q, validatorId)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		ve := structs.ValidatorEvent{}
		err = rows.Scan(&ve.ID, &ve.CreatedAt, &ve.UpdatedAt, &ve.ValidatorId, &ve.EventName, &ve.EventTime)
		if err != nil {
			return nil, err
		}
		validatorEvents = append(validatorEvents, ve)
	}
	if len(validatorEvents) == 0 {
		return nil, ErrNotFound
	}
	return validatorEvents, nil
}

// GetAllValidatorEvents gets all validator events
func (d *Driver) GetAllValidatorEvents(ctx context.Context) (validatorEvents []structs.ValidatorEvent, err error) {
	q := fmt.Sprintf("%s%s", getByStatementForValidatorEvent, orderByEventTime)
	rows, err := d.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		ve := structs.ValidatorEvent{}
		err = rows.Scan(&ve.ID, &ve.CreatedAt, &ve.UpdatedAt, &ve.ValidatorId, &ve.EventName, &ve.EventTime)
		if err != nil {
			return nil, err
		}
		validatorEvents = append(validatorEvents, ve)
	}
	if len(validatorEvents) == 0 {
		return nil, ErrNotFound
	}
	return validatorEvents, nil
}

package postgres

import (
	"../../structs"
	"context"
	"database/sql"
	"fmt"
)

const (
	insertStatementForValidator    = `INSERT INTO validators ("created_at", "updated_at", "name", "validator_address", "requested_address", "description", "fee_rate", "registration_time",  "minimum_delegation_amount", "accept_new_requests" ) VALUES ( NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8) `
	updateStatementForValidator    = `UPDATE validators SET updated_at = NOW(), name = $1, validator_address = $2, requested_address = $3, description = $4, fee_rate = $5, registration_time = $6, minimum_delegation_amount = $7, accept_new_requests = $8  WHERE id = $9 `
	getByStatementForValidator     = `SELECT id, created_at, updated_at, name, validator_address, requested_address, description, fee_rate, registration_time, minimum_delegation_amount, accept_new_requests FROM validators WHERE `
	byIdForValidator               = "id =  $1 "
	byValidatorAddressForValidator = "validator_address =  $1 "
	byRequestedAddressForValidator = "requested_address =  $1 "
)

// SaveOrUpdateValidator saves or updates validator
func (d *Driver) SaveOrUpdateValidator(ctx context.Context, v structs.Validator) error {
	_, err := d.GetValidatorById(ctx, v.ID)
	if err != nil {
		_, err = d.db.Exec(insertStatementForValidator, v.Name, v.ValidatorAddress, v.RequestedAddress, v.Description, v.FeeRate, v.RegistrationTime, v.MinimumDelegationAmount, v.AcceptNewRequests)
	} else {
		_, err = d.db.Exec(updateStatementForValidator, v.Name, v.ValidatorAddress, v.RequestedAddress, v.Description, v.FeeRate, v.RegistrationTime, v.MinimumDelegationAmount, v.AcceptNewRequests, v.ID)
	}
	return nil
}

// SaveOrUpdateValidators saves or updates validators
func (d *Driver) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	for _, v := range validators {
		if err := d.SaveOrUpdateValidator(ctx, v); err != nil {
			return err
		}
	}
	return nil
}

// GetValidatorById gets validator by id
func (d *Driver) GetValidatorById(ctx context.Context, id *string) (res structs.Validator, err error) {
	dlg := structs.Validator{}
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byIdForValidator)

	row := d.db.QueryRowContext(ctx, q, id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", err)
	}

	err = row.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.Name, &dlg.ValidatorAddress, &dlg.RequestedAddress, &dlg.Description, &dlg.FeeRate, &dlg.RegistrationTime, &dlg.MinimumDelegationAmount, &dlg.AcceptNewRequests)
	if err == sql.ErrNoRows || !(*dlg.ID != "") {
		return res, ErrNotFound
	}
	return dlg, err
}

// GetValidatorByValidatorAddress gets validator by validator address
func (d *Driver) GetValidatorByValidatorAddress(ctx context.Context, validatorAddress *string) (res structs.Validator, err error) {
	dlg := structs.Validator{}
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byValidatorAddressForValidator)

	row := d.db.QueryRowContext(ctx, q, validatorAddress)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", err)
	}

	err = row.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.Name, &dlg.ValidatorAddress, &dlg.RequestedAddress, &dlg.Description, &dlg.FeeRate, &dlg.RegistrationTime, &dlg.MinimumDelegationAmount, &dlg.AcceptNewRequests)
	if err == sql.ErrNoRows || !(*dlg.ID != "") {
		return res, ErrNotFound
	}
	return dlg, err
}

// GetValidatorByRequestedAddress gets validator by request address
func (d *Driver) GetValidatorByRequestedAddress(ctx context.Context, requestAddress *string) (res structs.Validator, err error) {
	dlg := structs.Validator{}
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byRequestedAddressForValidator)

	row := d.db.QueryRowContext(ctx, q, requestAddress)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", err)
	}

	err = row.Scan(&dlg.ID, &dlg.CreatedAt, &dlg.UpdatedAt, &dlg.Name, &dlg.ValidatorAddress, &dlg.RequestedAddress, &dlg.Description, &dlg.FeeRate, &dlg.RegistrationTime, &dlg.MinimumDelegationAmount, &dlg.AcceptNewRequests)
	if err == sql.ErrNoRows || !(*dlg.ID != "") {
		return res, ErrNotFound
	}
	return dlg, err
}

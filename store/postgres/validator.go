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
	getByStatementForValidator     = `SELECT v.id, v.created_at, v.updated_at, v.name, v.validator_address, v.requested_address, v.description, v.fee_rate, v.registration_time, v.minimum_delegation_amount, v.accept_new_requests FROM validators v WHERE `
	byIdForValidator               = `v.id =  $1 `
	byValidatorAddressForValidator = `v.validator_address =  $1 `
	byRequestedAddressForValidator = `v.requested_address =  $1 `
)

func (d *Driver) saveOrUpdateValidator(ctx context.Context, v structs.Validator) error {
	var err error
	if v.ID == "" {
		_, err = d.db.Exec(insertStatementForValidator, v.Name, v.ValidatorAddress, v.RequestedAddress, v.Description, v.FeeRate, v.RegistrationTime, v.MinimumDelegationAmount, v.AcceptNewRequests)
	} else {
		_, err = d.db.Exec(updateStatementForValidator, v.Name, v.ValidatorAddress, v.RequestedAddress, v.Description, v.FeeRate, v.RegistrationTime, v.MinimumDelegationAmount, v.AcceptNewRequests, v.ID)
	}
	return err
}

// SaveOrUpdateValidators saves or updates validators
func (d *Driver) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	for _, v := range validators {
		if err := d.saveOrUpdateValidator(ctx, v); err != nil {
			return err
		}
	}
	return nil
}

// GetValidatorById gets validator by id
func (d *Driver) GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error) {
	vld := structs.Validator{}
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byIdForValidator)

	row := d.db.QueryRowContext(ctx, q, id)
	if row.Err() != nil {
		return res, fmt.Errorf("query error: %w", row.Err().Error())
	}

	err = row.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.Name, &vld.ValidatorAddress, &vld.RequestedAddress, &vld.Description, &vld.FeeRate, &vld.RegistrationTime, &vld.MinimumDelegationAmount, &vld.AcceptNewRequests)
	if err == sql.ErrNoRows || !(vld.ID != "") {
		return res, ErrNotFound
	}
	return vld, err
}

// GetValidatorsByValidatorAddress gets validators by validator address
func (d *Driver) GetValidatorsByValidatorAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error) {
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byValidatorAddressForValidator)
	rows, err := d.db.QueryContext(ctx, q, validatorAddress)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.Name, &vld.ValidatorAddress, &vld.RequestedAddress, &vld.Description, &vld.FeeRate, &vld.RegistrationTime, &vld.MinimumDelegationAmount, &vld.AcceptNewRequests)
		if err != nil {
			return nil, err
		}
		validators = append(validators, vld)
	}
	if len(validators) == 0 {
		return nil, ErrNotFound
	}
	return validators, nil
}

// GetValidatorsByRequestedAddress gets validators by request address
func (d *Driver) GetValidatorsByRequestedAddress(ctx context.Context, requestAddress string) (validators []structs.Validator, err error) {
	q := fmt.Sprintf("%s%s", getByStatementForValidator, byRequestedAddressForValidator)
	rows, err := d.db.QueryContext(ctx, q, requestAddress)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vld := structs.Validator{}
		err = rows.Scan(&vld.ID, &vld.CreatedAt, &vld.UpdatedAt, &vld.Name, &vld.ValidatorAddress, &vld.RequestedAddress, &vld.Description, &vld.FeeRate, &vld.RegistrationTime, &vld.MinimumDelegationAmount, &vld.AcceptNewRequests)
		if err != nil {
			return nil, err
		}
		validators = append(validators, vld)
	}
	if len(validators) == 0 {
		return nil, ErrNotFound
	}
	return validators, nil

}

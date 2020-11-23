package handler

import (
	"github.com/figment-networks/skale-indexer/structs"
)

func validateDelegationRequiredFields(delegation structs.Delegation) error {
	if delegation.Holder == "" || delegation.ValidatorId == 0 ||
		delegation.Created.IsZero() ||
		delegation.Started.IsZero() || delegation.Finished.IsZero() ||
		delegation.Info == "" {
		return ErrMissingParameter
	}
	return nil
}

func validateDelegationEventRequiredFields(dlg structs.DelegationEvent) error {
	if dlg.DelegationId == "" || dlg.EventName == "" || dlg.EventTime.IsZero() {
		return ErrMissingParameter
	}
	return nil
}

func validateValidatorRequiredFields(validator structs.Validator) error {
	if validator.Name == "" || validator.Address == "" || validator.RequestedAddress == "" ||
		validator.Description == "" || validator.FeeRate == 0 ||
		validator.RegistrationTime.IsZero() {
		return ErrMissingParameter
	}
	return nil
}

func validateValidatorEventRequiredFields(ve structs.ValidatorEvent) error {
	if ve.ValidatorId == "" || ve.EventName == "" || ve.EventTime.IsZero() {
		return ErrMissingParameter
	}
	return nil
}

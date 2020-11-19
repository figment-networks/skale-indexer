package handler

import (
	"../structs"
)

func validateDelegationRequiredFields(delegation structs.Delegation) error {
	if delegation.Holder == nil || delegation.ValidatorId == nil || delegation.Amount == nil ||
		delegation.DelegationPeriod == nil || delegation.Created == nil ||
		delegation.Started == nil || delegation.Finished == nil ||
		delegation.Info == nil {
		return ErrMissingParameter
	}
	return nil
}

func validateDelegationEventRequiredFields(dlg structs.DelegationEvent) error {
	if dlg.DelegationId == nil || dlg.EventName == nil || dlg.EventTime == nil {
		return ErrMissingParameter
	}
	return nil
}

func validateValidatorRequiredFields(validator structs.Validator) error {
	if validator.Name == nil || validator.ValidatorAddress == nil || validator.RequestedAddress == nil ||
		validator.Description == nil || validator.FeeRate == nil ||
		validator.RegistrationTime == nil || validator.MinimumDelegationAmount == nil ||
		validator.AcceptNewRequests == nil {
		return ErrMissingParameter
	}
	return nil
}

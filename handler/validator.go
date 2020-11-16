package handler

import (
	"../structs"
)

func validateRequiredFields(delegation structs.Delegation) error {
	if delegation.Holder == nil || delegation.ValidatorId == nil || delegation.Amount == nil ||
		delegation.DelegationPeriod == nil || delegation.Created == nil ||
		delegation.Started == nil || delegation.Finished == nil ||
		delegation.Info == nil {
		return ErrMissingParameter
	}
	return nil
}

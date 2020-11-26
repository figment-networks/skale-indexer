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

func validateEventRequiredFields(dlg structs.Event) error {
	if dlg.BlockHeight == 0 || dlg.SmartContractAddress == "" ||
		dlg.TransactionIndex == 0 || dlg.EventType == "" ||
		dlg.EventName == "" || dlg.EventTime.IsZero() {
		return ErrMissingParameter
	}
	return nil
}

func validateValidatorRequiredFields(validator structs.Validator) error {
	if validator.Name == "" || len(validator.Address) == 0 ||
		validator.Description == "" || validator.FeeRate == 0 {
		return ErrMissingParameter
	}
	return nil
}

func validateNodeRequiredFields(node structs.Node) error {
	if node.Name == "" || node.Ip == "" || node.PublicIp == "" || node.Port == 0 ||
		node.PublicKey == "" || node.StartBlock == 0 || node.LastRewardDate.IsZero() ||
		node.FinishTime.IsZero() || node.Status == "" || node.ValidatorId == 0 {
		return ErrMissingParameter
	}
	return nil
}

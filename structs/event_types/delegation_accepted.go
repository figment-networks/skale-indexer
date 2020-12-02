package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/DelegationController.sol#L145
type EventDelegationAccepted struct {
	DelegationId uint64
}

func (e *EventDelegationAccepted) Trigger() error {
	/*
		delegation <- get delegation from 'delegations' table by delegation_id
	 	validator  <- get validator from 'validators' table by delegation.validator_id
		stats <- get stats from 'delegation_stats' table by validators.validator_id and status=ACCEPTED and statistics type = states
		increment stats 1
	*/
	return nil
}

package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/DelegationController.sol#L145
type EventDelegationAccepted struct {
	DelegationId uint64
}

func (e *EventDelegationAccepted) Trigger() error {
	return nil
}

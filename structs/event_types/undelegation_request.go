package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/DelegationController.sol#L159
type EventUnDelegationRequested struct {
	DelegationId uint64
}

func (e *EventUnDelegationRequested) Trigger() error {
	return nil
}

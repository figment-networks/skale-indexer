package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/ValidatorService.sol#L93
type EventNodeAddressWasAdded struct {
	ValidatorId uint64
	NodeAddress uint64
}

func (e *EventNodeAddressWasAdded) Trigger() error {
	return nil
}

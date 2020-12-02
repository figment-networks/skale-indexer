package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/Punisher.sol#L42
type EventSlash struct {
	ValidatorId uint64
	Amount      uint64
}

func (e *EventSlash) Trigger() error {
	return nil
}

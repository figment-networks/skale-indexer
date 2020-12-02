package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/Distributor.sol#L68
type EventBountyWasPaid struct {
	ValidatorId uint64
	Amount      uint64
}

func (e *EventBountyWasPaid) Trigger() error {
	return nil
}

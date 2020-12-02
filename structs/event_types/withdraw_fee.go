package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/Distributor.sol#L59
type EventUnWithdrawFee struct {
	ValidatorId uint64
	Destination uint64
	Amount      uint64
}

func (e *EventUnWithdrawFee) Trigger() error {
	return nil
}

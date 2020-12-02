package event_types

// see: https://github.com/skalenetwork/skale-manager/blob/24dc9c951674f2a217ecf3dafa1493d74758d15b/contracts/delegation/ValidatorService.sol#L71
type EventValidatorAddressChanged struct {
	ValidatorId uint64
	NewAddress  uint64
}

func (e *EventValidatorAddressChanged) Trigger() error {
	return nil
}
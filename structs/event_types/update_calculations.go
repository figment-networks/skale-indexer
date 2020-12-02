package event_types

type CalculationUpdater interface {
	Trigger() error
}

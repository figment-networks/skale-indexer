package contractor

type ClientContractor interface {
	delegationContractor
	delegationEventContractor
	validatorContractor
	validatorEventStore
}

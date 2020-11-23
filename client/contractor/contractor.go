package contractor

type ClientContractor interface {
	delegationContractor
	eventContractor
	validatorContractor
}

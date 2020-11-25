package structs

import "time"

type QueryParams struct {
	Id          string
	ValidatorId uint64
	Holder      string
	Address     []Address
	TimeFrom    time.Time
	TimeTo      time.Time
}

package structs

import "time"

type QueryParams struct {
	Id       string
	Address  Address
	TimeFrom time.Time
	TimeTo   time.Time
}

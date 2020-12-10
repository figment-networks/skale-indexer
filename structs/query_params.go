package structs

import "time"

type QueryParams struct {
	Id       string
	TimeFrom time.Time
	TimeTo   time.Time
}

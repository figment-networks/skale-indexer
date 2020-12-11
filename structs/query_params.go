package structs

import (
	"time"
)

type QueryParams struct {
	Id          string
	ValidatorId uint64
	TimeFrom    time.Time
	TimeTo      time.Time
}

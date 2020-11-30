package structs

import "time"

type QueryParams struct {
	Id            string
	ValidatorId   uint64
	Holder        uint64
	Address       []Address
	StatisticType StatisticType
	TimeFrom      time.Time
	TimeTo        time.Time
}

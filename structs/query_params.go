package structs

import "time"

type QueryParams struct {
	Id              string
	ValidatorId     uint64
	Holder          uint64
	Address         []Address
	ETHBlockHeight  uint64
	StatisticTypeDS StatisticTypeDS
	StatisticTypeVS StatisticTypeVS
	Status          uint64
	TimeFrom        time.Time
	TimeTo          time.Time
}

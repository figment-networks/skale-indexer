package structs

import (
	"time"
)

type Node struct {
	ID                       string    `json:"id"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
	Address                  uint64    `json:"address"`
	Name                     string    `json:"name"`
	Ip                       string    `json:"ip"`
	PublicIp                 string    `json:"public_ip"`
	Port                     uint      `json:"port"`
	PublicKey                string    `json:"public_key"`
	StartBlock               uint64    `json:"start_block"`
	LastRewardDate           time.Time `json:"last_reward_date"`
	FinishTime               time.Time `json:"finish_time"`
	Status                   string    `json:"status"`
	ValidatorId              uint64    `json:"validator_id"`
	RegistrationDate         time.Time `json:"registration_date"`
	LastBountyCall           time.Time `json:"last_bounty_call"`
	CalledGetBountyThisEpoch bool      `json:"called_get_bounty_this_epoch"`
	Balance                  float64   `json:"balance"`
}

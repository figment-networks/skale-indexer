package structs

import (
	"math/big"
	"net"
	"time"
)

type Node struct {
	ID             string     `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	NodeID         *big.Int   `json:"node_id"`
	Name           string     `json:"name"`
	IP             net.IP     `json:"ip"`
	PublicIP       net.IP     `json:"public_ip"`
	Port           uint16     `json:"port"`
	StartBlock     *big.Int   `json:"start_block"`
	NextRewardDate time.Time  `json:"next_reward_date"`
	LastRewardDate time.Time  `json:"last_reward_date"`
	FinishTime     *big.Int   `json:"finish_time"`
	Status         NodeStatus `json:"node_status"`
	ValidatorID    *big.Int   `json:"validator_id"`
	EventTime      time.Time  `json:"event_time"`
	BlockHeight    uint64     `json:"block_height"`
}

type NodeStatus uint

const (
	NodeStatusActive NodeStatus = iota
	NodeStatusLeaving
	NodeStatusLeft
	NodeStatusInMaintenance
)

func (k NodeStatus) String() string {
	switch k {
	case NodeStatusActive:
		return "Active"
	case NodeStatusLeaving:
		return "Leaving"
	case NodeStatusLeft:
		return "Left"
	case NodeStatusInMaintenance:
		return "In_Maintenance"
	default:
		return "unknown"
	}
}

package postgresql

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// SaveNode saves node
func (d *Driver) SaveNode(ctx context.Context, n structs.Node) error {
	_, err := d.db.Exec(`INSERT INTO nodes
			("node_id", "name",  "ip", "public_ip", "port", "start_block", "next_reward_date", "last_reward_date", "finish_time", "status", "validator_id", "event_time")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (node_id)
			DO UPDATE SET
				name = EXCLUDED.name,
				ip = EXCLUDED.ip,
				public_ip = EXCLUDED.public_ip,
				port = EXCLUDED.port,
				start_block = EXCLUDED.start_block,
				next_reward_date = EXCLUDED.next_reward_date,
				last_reward_date = EXCLUDED.last_reward_date,
				finish_time = EXCLUDED.finish_time,
				status = EXCLUDED.status,
				validator_id = EXCLUDED.validator_id,
				event_time = EXCLUDED.event_time
			`,
		n.NodeID.String(),
		n.Name,
		n.IP.String(),
		n.PublicIP.String(),
		n.Port,
		n.StartBlock.String(),
		n.NextRewardDate,
		n.LastRewardDate,
		n.FinishTime.String(),
		n.Status,
		n.ValidatorID.String(),
		n.EventTime)
	return err
}

// GetNodes gets nodes
func (d *Driver) GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error) {
	q := `SELECT
			id, created_at, node_id, name, ip, public_ip, port, start_block, next_reward_date, last_reward_date, finish_time, status, validator_id, event_time
		FROM nodes `

	var (
		args   []interface{}
		wherec []string
		i      = 1
	)

	if params.NodeId != "" {
		wherec = append(wherec, ` node_id =  $`+strconv.Itoa(i))
		args = append(args, params.NodeId)
		i++
	}
	if params.ValidatorId != "" {
		wherec = append(wherec, ` validator_id =  $`+strconv.Itoa(i))
		args = append(args, params.ValidatorId)
		i++
	}
	if len(args) > 0 {
		q += ` WHERE `
		q += strings.Join(wherec, " AND ")
	}

	q += ` ORDER BY node_id, event_time`

	rows, err := d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		n := structs.Node{}
		var nodeId uint64
		var startBlock uint64
		var finishTime uint64
		var validatorId uint64
		var IP string
		var publicIP string

		err = rows.Scan(&n.ID, &n.CreatedAt, &nodeId, &n.Name, &IP, &publicIP, &n.Port, &startBlock, &n.NextRewardDate, &n.LastRewardDate, &finishTime, &n.Status, &validatorId, &n.EventTime)
		if err != nil {
			return nil, err
		}

		n.NodeID = new(big.Int).SetUint64(nodeId)
		n.StartBlock = new(big.Int).SetUint64(startBlock)
		n.FinishTime = new(big.Int).SetUint64(finishTime)
		n.ValidatorID = new(big.Int).SetUint64(validatorId)
		n.IP, _, _ = net.ParseCIDR(IP)
		n.PublicIP, _, _ = net.ParseCIDR(publicIP)
		nodes = append(nodes, n)
	}
	return nodes, nil
}

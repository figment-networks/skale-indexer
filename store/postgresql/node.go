package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"net"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// SaveNode saves node
func (d *Driver) SaveNode(ctx context.Context, n structs.Node) error {
	_, err := d.db.Exec(`INSERT INTO nodes
			("node_id", "name",  "ip", "public_ip", "port", "start_block", "next_reward_date", "last_reward_date", "finish_time", "status", "validator_id", "event_time")
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
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
	var q string
	var rows *sql.Rows

	q = `SELECT
			id, created_at, node_id, name, ip, public_ip, port, start_block, next_reward_date, last_reward_date, finish_time, status, validator_id, event_time
		FROM nodes `

	if params.ValidatorId != "" {
		q += ` WHERE validator_id =  $1 `
	}
	q += ` ORDER BY node_id, event_time`

	//byRecentStartBlockN = `AND start_block = (SELECT n2.start_block FROM nodes n2 WHERE n2.validator_id = $2 ORDER BY n2.start_block DESC LIMIT 1) `

	//	if params.ValidatorId != "" && !params.Recent {
	//	q = fmt.Sprintf("%s%s%s", getByStatementN, byValidatorIdN, orderByNameN)
	//		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	//	} else if params.ValidatorId != "" && params.Recent {
	//	q = fmt.Sprintf("%s%s%s%s", getByStatementN, byValidatorIdN, byRecentStartBlockN, orderByNameN)
	//		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId, params.ValidatorId)
	//	} else {
	//		q = fmt.Sprintf("%s%s", getByStatementN, orderByNameN)
	//		rows, err = d.db.QueryContext(ctx, q)
	//	}

	rows, err = d.db.QueryContext(ctx, q)
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

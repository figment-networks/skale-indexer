package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"net"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// TODO: run explain analyze to check full scan and add required indexes
const (
	getByStatementN     = `SELECT id, created_at, node_id, name, ip, public_ip, port, start_block, next_reward_date, last_reward_date, finish_time, status, validator_id, event_time FROM nodes `
	byValidatorIdN      = `WHERE validator_id =  $1 `
	byRecentStartBlockN = `AND start_block =  (SELECT n2.start_block FROM nodes n2 WHERE n2.validator_id = $2 ORDER BY n2.start_block DESC LIMIT 1) `
	orderByNameN        = `ORDER BY name `
)

// SaveNode saves node
func (d *Driver) SaveNode(ctx context.Context, n structs.Node) error {
	_, err := d.db.Exec(`INSERT INTO nodes ("node_id", "name", "ip", "public_ip", "port", "start_block", "next_reward_date", "last_reward_date", "finish_time", "status", "validator_id", "event_time") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		n.NodeID.String(),
		n.Name,
		net.IPv4(n.IP[0], n.IP[1], n.IP[2], n.IP[3]),
		net.IPv4(n.PublicIP[0], n.PublicIP[1], n.PublicIP[2], n.PublicIP[3]),
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
	if params.ValidatorId != "" && !params.Recent {
		q = fmt.Sprintf("%s%s%s", getByStatementN, byValidatorIdN, orderByNameN)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else if params.ValidatorId != "" && params.Recent {
		q = fmt.Sprintf("%s%s%s%s", getByStatementN, byValidatorIdN, byRecentStartBlockN, orderByNameN)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId, params.ValidatorId)
	} else {
		q = fmt.Sprintf("%s%s", getByStatementN, orderByNameN)
		rows, err = d.db.QueryContext(ctx, q)
	}

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
		var ip string
		var publicIp string

		err = rows.Scan(&n.ID, &n.CreatedAt, &nodeId, &n.Name, &ip, &publicIp, &n.Port, &startBlock, &n.NextRewardDate, &n.LastRewardDate, &finishTime, &n.Status, &validatorId, &n.EventTime)
		if err != nil {
			return nil, err
		}

		n.NodeID = new(big.Int).SetUint64(nodeId)
		n.StartBlock = new(big.Int).SetUint64(startBlock)
		n.FinishTime = new(big.Int).SetUint64(finishTime)
		n.ValidatorID = new(big.Int).SetUint64(validatorId)

		k := strings.Split(ip, ".")
		var ipp [4]byte
		for i := 0; i < len(ipp); i++ {
			z := []byte(k[i])
			ipp[i] = z[0]
		}
		n.IP = ipp
		k = strings.Split(publicIp, ".")
		var pip [4]byte
		for i := 0; i < len(pip); i++ {
			z := []byte(k[i])
			pip[i] = z[0]
		}
		n.PublicIP = pip

		nodes = append(nodes, n)
	}
	if len(nodes) == 0 {
		return nil, ErrNotFound
	}
	return nodes, nil
}

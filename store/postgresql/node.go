package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/client/structs"
	"github.com/figment-networks/skale-indexer/handler"
)

const (
	insertStatementN        = `INSERT INTO nodes ("node_id", "name", "ip", "public_ip", "port", "start_block", "next_reward_date", "last_reward_date", "finish_time", "status", "validator_id", "event_time") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `
	getByStatementN         = `SELECT n.id, n.created_at, n.updated_at, n.node_id, n.name, n.ip, n.public_ip, n.port, n.start_block, n.next_reward_date, n.last_reward_date, n.finish_time, n.status, n.validator_id, n.event_time FROM nodes n `
	byIdN                   = `WHERE n.id =  $1 `
	byValidatorIdN          = `WHERE n.validator_id =  $1 `
	byRecentStartBlockN     = `AND n.start_block =  (SELECT n2.start_block FROM nodes n2 WHERE n2.validator_id = $2 ORDER BY n2.start_block DESC LIMIT 1) `
	orderByNameN            = `ORDER BY n.name DESC `
)

// SaveNode saves node
func (d *Driver) SaveNode(ctx context.Context, n structs.Node) error {
	_, err := d.db.Exec(insertStatementN, n.NodeID, n.Name, n.IP, n.PublicIP, n.Port, n.StartBlock, n.NextRewardDate, n.LastRewardDate, n.FinishTime, n.Status, n.ValidatorID, n.EventTime)
	return err
}

// GetNodes gets nodes
func (d *Driver) GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementN, byIdN)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else if params.ValidatorId != 0 && !params.Recent {
		q = fmt.Sprintf("%s%s%s", getByStatementN, byValidatorIdN, orderByNameN)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else if params.ValidatorId != 0 && params.Recent {
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
		err = rows.Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt, &n.NodeID, &n.Name, &n.IP, &n.PublicIP, &n.Port, &n.StartBlock, &n.NextRewardDate, &n.LastRewardDate, &n.FinishTime, &n.Status, &n.ValidatorID, &n.EventTime)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	if len(nodes) == 0 {
		return nil, handler.ErrNotFound
	}
	return nodes, nil
}

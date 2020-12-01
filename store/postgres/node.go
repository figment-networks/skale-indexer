package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	insertStatementForNode = `INSERT INTO nodes ("updated_at", "address", "name", "ip", "public_ip", "port", "public_key", "start_block", "last_reward_date", "finish_time", "status", "validator_id", "registration_date", "last_bounty_call", "called_get_bounty_this_epoch", "balance") VALUES ( NOW(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) `
	updateStatementForNode = `UPDATE nodes SET updated_at = NOW(), address = $1,  name = $2, ip = $3, public_ip = $4, port = $5, public_key = $6, start_block = $7, last_reward_date = $8, finish_time = $9, status = $10, validator_id = $11, registration_date = $12, last_bounty_call = $13, called_get_bounty_this_epoch = $14, balance = $15 WHERE id = $16 `
	getByStatementForNode  = `SELECT n.id, n.created_at, n.updated_at, n.address, n.name, n.ip, n.public_ip, n.port, n.public_key, n.start_block, n.last_reward_date, n.finish_time, n.status, n.validator_id, n.registration_date, n.last_bounty_call, n.called_get_bounty_this_epoch, n.balance FROM nodes n `
	byIdForNode            = `WHERE n.id =  $1 `
	byValidatorIdForNode   = `WHERE n.validator_id =  $1 `
	orderByName            = `ORDER BY n.name DESC`
)

func (d *Driver) saveOrUpdateNode(ctx context.Context, n structs.Node) error {
	if n.ID == "" {
		_, err := d.db.Exec(insertStatementForNode, n.Address, n.Name, n.Ip, n.PublicIp, n.Port, n.PublicKey, n.StartBlock, n.LastRewardDate, n.FinishTime, n.Status, n.ValidatorId, n.RegistrationDate, n.LastBountyCall, n.CalledGetBountyThisEpoch, n.Balance)
		return err
	}
	_, err := d.db.Exec(updateStatementForNode, n.Address, n.Name, n.Ip, n.PublicIp, n.Port, n.PublicKey, n.StartBlock, n.LastRewardDate, n.FinishTime, n.Status, n.ValidatorId, n.RegistrationDate, n.LastBountyCall, n.CalledGetBountyThisEpoch, n.Balance, n.ID)
	return err
}

// SaveOrUpdateNodes saves or updates nodes
func (d *Driver) SaveOrUpdateNodes(ctx context.Context, nodes []structs.Node) error {
	for _, n := range nodes {
		if err := d.saveOrUpdateNode(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

// GetNodes gets nodes
func (d *Driver) GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForNode, byIdForNode)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else if params.ValidatorId > 0 {
		q = fmt.Sprintf("%s%s", getByStatementForNode, byValidatorIdForNode)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else {
		q = fmt.Sprintf("%s%s", getByStatementForNode, orderByName)
		rows, err = d.db.QueryContext(ctx, q)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		n := structs.Node{}
		err = rows.Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt, &n.Address, &n.Name, &n.Ip, &n.PublicIp, &n.Port, &n.PublicKey, &n.StartBlock, &n.LastRewardDate, &n.FinishTime, &n.Status, &n.ValidatorId, &n.RegistrationDate, &n.LastBountyCall, &n.CalledGetBountyThisEpoch, &n.Balance)
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

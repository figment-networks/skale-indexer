package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"math/big"
	"net"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

var zero common.Address

// SaveNodes saves nodes
func (d *Driver) SaveNodes(ctx context.Context, nodes []structs.Node, removedNodeAddress common.Address) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// upsert nodes
	for _, n := range nodes {
		params := structs.NodeParams{
			NodeID: n.NodeID.String(),
		}
		no, err := d.GetNodes(ctx, params)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		if no == nil || (len(no) > 0 && no[0].BlockHeight <= n.BlockHeight) {
			_, err = tx.ExecContext(ctx, `INSERT INTO nodes
			("node_id", "address", "name",  "ip", "public_ip", "port", "start_block", "next_reward_date", "last_reward_date", "finish_time", "status", "validator_id", "block_height")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			ON CONFLICT (node_id)
			DO UPDATE SET
				name = EXCLUDED.name,
				address = EXCLUDED.address,
				ip = EXCLUDED.ip,
				public_ip = EXCLUDED.public_ip,
				port = EXCLUDED.port,
				start_block = EXCLUDED.start_block,
				next_reward_date = EXCLUDED.next_reward_date,
				last_reward_date = EXCLUDED.last_reward_date,
				finish_time = EXCLUDED.finish_time,
				status = EXCLUDED.status,
				validator_id = EXCLUDED.validator_id,
				block_height = EXCLUDED.block_height`,
				n.NodeID.String(),
				n.Address.Hash().Big().String(),
				n.Name,
				n.IP.String(),
				n.PublicIP.String(),
				n.Port,
				n.StartBlock.String(),
				n.NextRewardDate,
				n.LastRewardDate,
				n.FinishTime.String(),
				n.Status.String(),
				n.ValidatorID.String(),
				n.BlockHeight)
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	// update removed node
	if removedNodeAddress.Hash().Big().String() != zero.Hash().Big().String() && len(nodes) > 0 {
		params := structs.NodeParams{
			Address:     removedNodeAddress.Hash().Big().String(),
			ValidatorID: nodes[0].ValidatorID.String(),
		}
		no, err := d.GetNodes(ctx, params)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		if no == nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return errors.New("removed node address is not available on db")
		}

		if no[0].BlockHeight <= nodes[0].BlockHeight {
			nodeIds := make([]int64, len(nodes))
			for i, n := range nodes {
				nodeIds[i] = n.NodeID.Int64()
			}
			_, err = tx.ExecContext(ctx, `UPDATE nodes SET address = $1 
				WHERE validator_id = $2 AND address = $3 AND node_id NOT IN ($4)`,
				zero.Hash().Big().String(),
				nodes[0].ValidatorID.Int64(),
				removedNodeAddress.Hash().Big().String(),
				pq.Array(nodeIds),
			)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return rollbackErr
				}
				return err
			}
		}
	}

	return tx.Commit()
}

// GetNodes gets nodes
func (d *Driver) GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error) {
	q := `SELECT
			id, created_at, node_id, address, name, ip, public_ip, port, start_block, next_reward_date, last_reward_date, finish_time, status, validator_id, block_height
		FROM nodes `

	var (
		args   []interface{}
		wherec []string
		i      = 1
	)

	if params.NodeID != "" {
		wherec = append(wherec, ` node_id =  $`+strconv.Itoa(i))
		args = append(args, params.NodeID)
		i++
	}
	if params.ValidatorID != "" {
		wherec = append(wherec, ` validator_id =  $`+strconv.Itoa(i))
		args = append(args, params.ValidatorID)
		i++
	}
	if params.Status != "" {
		wherec = append(wherec, ` status =  $`+strconv.Itoa(i))
		args = append(args, params.Status)
		i++
	}
	if params.Address != "" {
		wherec = append(wherec, ` address =  $`+strconv.Itoa(i))
		args = append(args, common.HexToAddress(params.Address).Hash().Big().String())
		i++
	}
	if len(args) > 0 {
		q += ` WHERE `
		q += strings.Join(wherec, " AND ")
	}

	q += ` ORDER BY node_id, start_block`

	rows, err := d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		n := structs.Node{}
		var nodeId uint64
		var address []byte
		var startBlock uint64
		var finishTime uint64
		var validatorId uint64
		var IP string
		var publicIP string
		var status string
		err = rows.Scan(&n.ID, &n.CreatedAt, &nodeId, &address, &n.Name, &IP, &publicIP, &n.Port, &startBlock, &n.NextRewardDate, &n.LastRewardDate, &finishTime, &status, &validatorId, &n.BlockHeight)
		if err != nil {
			return nil, err
		}

		n.NodeID = new(big.Int).SetUint64(nodeId)
		a := new(big.Int)
		a.SetString(string(address), 10)
		n.Address.SetBytes(a.Bytes())
		n.StartBlock = new(big.Int).SetUint64(startBlock)
		n.FinishTime = new(big.Int).SetUint64(finishTime)
		n.ValidatorID = new(big.Int).SetUint64(validatorId)
		n.IP, _, _ = net.ParseCIDR(IP)
		n.PublicIP, _, _ = net.ParseCIDR(publicIP)
		s, _ := structs.GetTypeForNode(status)
		n.Status = s
		nodes = append(nodes, n)
	}
	return nodes, nil
}

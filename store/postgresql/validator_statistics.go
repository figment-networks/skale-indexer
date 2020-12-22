package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// TODO: run explain analyze to check full scan and add required indexes
const (
	calculateTotalStake = `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
 								SELECT t1.validator_id, sum(t1.amount) AS amount, $1 AS block_height, $2 AS statistics_type FROM
 									(SELECT  DISTINCT ON (delegation_id) validator_id, delegation_id, block_height, state, amount FROM delegations
 											WHERE validator_id = $3 AND block_height <=$4 ORDER BY delegation_id, block_height DESC)  t1
 								WHERE  t1.state IN ($5, $6) GROUP BY t1.validator_id
 							ON CONFLICT (statistics_type, validator_id, block_height)
 							DO UPDATE SET amount = EXCLUDED.amount`
	calculateActiveNodes = `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount,  $1 AS block_height, $2 AS statistics_type FROM nodes
									WHERE validator_id = $3 AND status = $4 GROUP BY validator_id LIMIT 1)
							ON CONFLICT (statistics_type, validator_id, block_height)
 							DO UPDATE SET amount = EXCLUDED.amount		`
	calculateLinkedNodes = `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount, $1 AS block_height, $2 AS statistics_type FROM nodes
									WHERE validator_id = $3  GROUP BY validator_id LIMIT 1)
							ON CONFLICT (statistics_type, validator_id, block_height)
 							DO UPDATE SET amount = EXCLUDED.amount`

	updateTotalStake = `UPDATE validators SET staked = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ORDER BY block_height DESC LIMIT 1) 
							WHERE validator_id = $4 `
	updateActiveNodes = `UPDATE validators SET active_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ORDER BY block_height DESC LIMIT 1)
 							WHERE validator_id = $4 `
	updateLinkedNodes = `UPDATE validators SET linked_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ORDER BY block_height DESC LIMIT 1)
 							WHERE validator_id = $4 `
)

func (d *Driver) GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	q := `SELECT id, created_at, validator_id, amount, block_height, statistics_type FROM validator_statistics WHERE statistics_type = $1 `
	var rows *sql.Rows
	if params.ValidatorId != "" {
		q = fmt.Sprintf("%s%s%s", q, `AND validator_id = $2 `, `ORDER BY block_height DESC `)
		if params.Recent {
			q = fmt.Sprintf("%s%s", q, `LIMIT 1`)
		}
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS, params.ValidatorId)
	} else {
		// return latest for each validator
		q = `SELECT DISTINCT ON (validator_id)  id, created_at, validator_id, amount, block_height, statistics_type FROM validator_statistics 
					WHERE statistics_type = $1 ORDER BY validator_id, block_height DESC `
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		vs := structs.ValidatorStatistics{}
		var vldId uint64
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &vldId, &vs.Amount, &vs.BlockHeight, &vs.StatisticType)
		vs.ValidatorId = new(big.Int).SetUint64(vldId)
		if err != nil {
			return nil, err
		}
		validatorStatistics = append(validatorStatistics, vs)
	}
	return validatorStatistics, nil
}

// update "TOTAL STAKE" = Accepted + UndelegationRequested
func (d *Driver) CalculateTotalStake(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	_, err := d.db.Exec(calculateTotalStake, params.BlockHeight, structs.ValidatorStatisticsTypeTotalStake, params.ValidatorId, params.BlockHeight, structs.DelegationStateACCEPTED, structs.DelegationStateUNDELEGATION_REQUESTED)
	if err == nil {
		_, err = d.db.Exec(updateTotalStake, params.ValidatorId, structs.ValidatorStatisticsTypeTotalStake, params.ValidatorId, params.BlockHeight, params.ValidatorId, params.BlockHeight)
	}
	return err
}

// update "ACTIVE NODES"
func (d *Driver) CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	_, err := d.db.Exec(calculateActiveNodes, params.BlockHeight, structs.ValidatorStatisticsTypeActiveNodes, params.ValidatorId, structs.NodeStatusActive)
	if err == nil {
		_, err = d.db.Exec(updateActiveNodes, params.ValidatorId, structs.ValidatorStatisticsTypeActiveNodes, params.BlockHeight, params.ValidatorId)
	}
	return err
}

// update "LINKED NODES"
func (d *Driver) CalculateLinkedNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	_, err := d.db.Exec(calculateLinkedNodes, params.BlockHeight, structs.ValidatorStatisticsTypeLinkedNodes, params.ValidatorId)
	if err == nil {
		_, err = d.db.Exec(updateLinkedNodes, params.ValidatorId, structs.ValidatorStatisticsTypeLinkedNodes, params.BlockHeight, params.ValidatorId)
	}
	return err
}

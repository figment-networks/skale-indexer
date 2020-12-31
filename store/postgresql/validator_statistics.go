package postgresql

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// SaveValidatorStatistics saves validator statistics
func (d *Driver) SaveValidatorStatistics(ctx context.Context, vs structs.ValidatorStatistics) error {
	_, err := d.db.Exec(`INSERT INTO validator_statistics ("validator_id", "amount", "block_height", "statistics_type") 
			VALUES ($1, $2, $3, $4) 
			ON CONFLICT (statistics_type, validator_id, block_height)
			DO UPDATE SET
			amount = EXCLUDED.amount `,
		vs.ValidatorId.String(),
		vs.Amount.String(),
		vs.BlockHeight,
		vs.StatisticType)
	return err
}

func (d *Driver) GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	q := `SELECT
			DISTINCT ON (validator_id, statistics_type) 
				id, created_at, validator_id, amount, block_height, statistics_type 
			FROM validator_statistics `
	var (
		args   []interface{}
		wherec []string
		i      = 1
	)
	if params.ValidatorId != "" {
		wherec = append(wherec, ` validator_id =  $`+strconv.Itoa(i))
		args = append(args, params.ValidatorId)
		i++
	}
	if params.StatisticsTypeVS != "" {
		wherec = append(wherec, ` statistics_type =  $`+strconv.Itoa(i))
		args = append(args, params.StatisticsTypeVS)
		i++
	}
	if len(args) > 0 {
		q += ` WHERE `
	}
	q += strings.Join(wherec, " AND ")
	q += ` ORDER BY validator_id ASC, statistics_type, block_height DESC`

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vs := structs.ValidatorStatistics{}
		var vldId uint64
		var amount string
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &vldId, &amount, &vs.BlockHeight, &vs.StatisticType)
		if err != nil {
			return nil, err
		}
		vs.ValidatorId = new(big.Int).SetUint64(vldId)
		amnt := new(big.Int)
		amnt.SetString(amount, 10)
		vs.Amount = amnt
		validatorStatistics = append(validatorStatistics, vs)
	}
	return validatorStatistics, nil
}

func (d *Driver) GetValidatorStatisticsChart(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	q := `SELECT id, created_at, validator_id, amount, block_height, statistics_type 
			FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 ORDER BY block_height DESC`
	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, params.ValidatorId, params.StatisticsTypeVS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vs := structs.ValidatorStatistics{}
		var vldId uint64
		var amount string
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &vldId, &amount, &vs.BlockHeight, &vs.StatisticType)
		if err != nil {
			return nil, err
		}
		vs.ValidatorId = new(big.Int).SetUint64(vldId)
		amnt := new(big.Int)
		amnt.SetString(amount, 10)
		vs.Amount = amnt
		validatorStatistics = append(validatorStatistics, vs)
	}
	return validatorStatistics, nil
}

func (d *Driver) CalculateTotalStake(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	calculateTotalStake := `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
 								SELECT t1.validator_id, sum(t1.amount) AS amount, $1 AS block_height, $2 AS statistics_type FROM
 									(SELECT  DISTINCT ON (delegation_id) validator_id, delegation_id, block_height, state, amount FROM delegations
 											WHERE validator_id = $3 AND block_height <=$4 ORDER BY delegation_id, block_height DESC)  t1
 								WHERE  t1.state IN ($5, $6) GROUP BY t1.validator_id
 							ON CONFLICT (statistics_type, validator_id, block_height)
 							DO UPDATE SET amount = EXCLUDED.amount`
	_, err := d.db.Exec(calculateTotalStake, params.BlockHeight, structs.ValidatorStatisticsTypeTotalStake, params.ValidatorId, params.BlockHeight, structs.DelegationStateACCEPTED, structs.DelegationStateUNDELEGATION_REQUESTED)
	if err == nil {
		updateTotalStake := `UPDATE validators SET staked = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ) WHERE validator_id = $4`
		_, err = d.db.Exec(updateTotalStake, params.ValidatorId, structs.ValidatorStatisticsTypeTotalStake, params.BlockHeight, params.ValidatorId)
	}
	return err
}

func (d *Driver) CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	calculateActiveNodes := `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount,  $1 AS block_height, $2 AS statistics_type FROM nodes
									WHERE validator_id = $3 AND status = $4 GROUP BY validator_id LIMIT 1)
							ON CONFLICT (statistics_type, validator_id, block_height)
 							DO UPDATE SET amount = EXCLUDED.amount `
	_, err := d.db.Exec(calculateActiveNodes, params.BlockHeight, structs.ValidatorStatisticsTypeActiveNodes, params.ValidatorId, structs.NodeStatusActive)
	if err == nil {
		updateActiveNodes := `UPDATE validators SET active_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ) WHERE validator_id = $4`
		_, err = d.db.Exec(updateActiveNodes, params.ValidatorId, structs.ValidatorStatisticsTypeActiveNodes, params.BlockHeight, params.ValidatorId)
	}
	return err
}

func (d *Driver) CalculateLinkedNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	calculateLinkedNodes := `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount, $1 AS block_height, $2 AS statistics_type FROM nodes
									WHERE validator_id = $3  GROUP BY validator_id LIMIT 1)
							ON CONFLICT (statistics_type, validator_id, block_height)
 							DO UPDATE SET amount = EXCLUDED.amount`
	_, err := d.db.Exec(calculateLinkedNodes, params.BlockHeight, structs.ValidatorStatisticsTypeLinkedNodes, params.ValidatorId)
	if err == nil {
		updateLinkedNodes := `UPDATE validators SET linked_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ) WHERE validator_id = $4`
		_, err = d.db.Exec(updateLinkedNodes, params.ValidatorId, structs.ValidatorStatisticsTypeLinkedNodes, params.BlockHeight, params.ValidatorId)
	}
	return err
}

func (d *Driver) UpdateUnclaimedRewards(ctx context.Context, validatorId *big.Int, blockHeight uint64) error {
	updateUnclaimedRewards := `UPDATE validators SET unclaimed_rewards = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 ) WHERE validator_id = $4`
	_, err := d.db.Exec(updateUnclaimedRewards, validatorId.String(), structs.ValidatorStatisticsTypeUnclaimedRewards, blockHeight, validatorId.String())
	return err
}

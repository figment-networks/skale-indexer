package postgresql

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

func (d *Driver) SaveValidatorStatistic(ctx context.Context, validatorID *big.Int, blockHeight uint64, blockTime time.Time, statisticsType structs.StatisticTypeVS, amount *big.Int) (err error) {
	// (lukanus): Update value in validator_statistics unless the value already exists
	_, err = d.db.ExecContext(ctx, `
	INSERT INTO validator_statistics (validator_id, block_height, time, statistic_type, amount)
		( SELECT $1, $2, $3, $4, $5	WHERE
			NOT EXISTS (
					SELECT 1 FROM validator_statistics
					WHERE validator_id = $1 AND statistic_type = $4 AND block_height < $2 AND amount = $5
					ORDER BY block_height DESC LIMIT 1
					)
		)
		ON CONFLICT (validator_id, block_height, statistic_type)
		DO UPDATE SET amount = EXCLUDED.amount;`,
		validatorID.String(), blockHeight, blockTime, statisticsType, amount.String())
	return nil
}

func (d *Driver) GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	q := `SELECT
			DISTINCT ON (validator_id, statistic_type) id, created_at, validator_id, amount, block_height, time, statistic_type
			FROM validator_statistics `
	var (
		args   []interface{}
		wherec []string
		i      = 1
	)
	if params.ValidatorID != "" {
		wherec = append(wherec, ` validator_id =  $`+strconv.Itoa(i))
		args = append(args, params.ValidatorID)
		i++
	}
	if params.Type > 0 {
		wherec = append(wherec, ` statistic_type =  $`+strconv.Itoa(i))
		args = append(args, params.Type)
		i++
	}
	if len(args) > 0 {
		q += ` WHERE `
	}
	q += strings.Join(wherec, " AND ")
	q += ` ORDER BY validator_id ASC, statistic_type, block_height DESC`

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		vldId  uint64
		amount string
	)

	for rows.Next() {
		vs := structs.ValidatorStatistics{}
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &vldId, &amount, &vs.BlockHeight, &vs.Time, &vs.Type)
		if err != nil {
			return nil, err
		}
		vs.ValidatorID = new(big.Int).SetUint64(vldId)
		amnt := new(big.Int)
		amnt.SetString(amount, 10)
		vs.Amount = amnt
		validatorStatistics = append(validatorStatistics, vs)
	}
	return validatorStatistics, nil
}

func (d *Driver) GetValidatorStatisticsTimeline(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx,
		`SELECT id, created_at, validator_id, amount, block_height, time, statistic_type
			FROM validator_statistics
			WHERE
				validator_id = $1 AND statistic_type = $2
			ORDER BY block_height DESC`, params.ValidatorID, params.Type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		vldId  uint64
		amount string
	)

	for rows.Next() {
		vs := structs.ValidatorStatistics{}
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &vldId, &amount, &vs.BlockHeight, &vs.Time, &vs.Type)
		if err != nil {
			return nil, err
		}
		vs.ValidatorID = new(big.Int).SetUint64(vldId)
		amnt := new(big.Int)
		amnt.SetString(amount, 10)
		vs.Amount = amnt
		validatorStatistics = append(validatorStatistics, vs)
	}
	return validatorStatistics, nil
}

func (d *Driver) CalculateTotalStake(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO validator_statistics (validator_id, block_height, time, statistic_type, amount)
										SELECT $1, $2, $3, $4, sum(t1.amount) AS amount
									FROM
											( SELECT DISTINCT ON (delegation_id) validator_id, delegation_id, block_height, state, amount
												FROM delegations
												WHERE validator_id = $1 AND block_height <=$2
												ORDER BY delegation_id, block_height DESC) t1
										WHERE  t1.state IN ($5, $6) GROUP BY t1.validator_id
									ON CONFLICT (validator_id, block_height, statistic_type)
									DO UPDATE SET amount = EXCLUDED.amount`,
		params.ValidatorID, params.BlockHeight, params.BlockTime, structs.ValidatorStatisticsTypeTotalStake, structs.DelegationStateACCEPTED, structs.DelegationStateUNDELEGATION_REQUESTED)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE validators
						SET staked = (
							 	SELECT COALESCE((SELECT amount
								 FROM validator_statistics
								 WHERE validator_id = $1 AND statistic_type = $2 
								 ORDER BY block_height DESC LIMIT 1 ), 0))
						WHERE validator_id = $1`,
		params.ValidatorID, structs.ValidatorStatisticsTypeTotalStake)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

func (d *Driver) CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	// BUG(l): This is wrong, would give random results. It either has to be calculated from state, or nodes

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO validator_statistics (validator_id, block_height, time, statistic_type, amount)
				(SELECT $1, $2, $3, $4, count(*) AS amount
					FROM nodes
					WHERE validator_id = $1 AND status = $5 AND address != $6
					GROUP BY validator_id LIMIT 1)
			ON CONFLICT (validator_id, block_height, statistic_type)
			DO UPDATE SET amount = EXCLUDED.amount `,
		params.ValidatorID, params.BlockHeight, params.BlockTime, structs.ValidatorStatisticsTypeActiveNodes, structs.NodeStatusActive.String(), zero)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE validators
						SET active_nodes =  (
							 	SELECT COALESCE((SELECT amount
								 FROM validator_statistics
								 WHERE validator_id = $1 AND statistic_type = $2 
								 ORDER BY block_height DESC LIMIT 1 ), 0))
						WHERE validator_id = $1`,
		params.ValidatorID, structs.ValidatorStatisticsTypeActiveNodes)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

func (d *Driver) CalculateLinkedNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	// BUG(l): This is wrong, would give random results. It either has to be calculated from  state, or nodes
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO validator_statistics (validator_id, block_height, time, statistic_type, amount)
									(SELECT  $1, $2, $3, $4, count(*) AS amount
									FROM nodes
									WHERE validator_id = $1 AND address != $5
									GROUP BY validator_id LIMIT 1)
								ON CONFLICT (validator_id, block_height, statistic_type)
								DO UPDATE SET amount = EXCLUDED.amount`,
		params.ValidatorID, params.BlockHeight, params.BlockTime, structs.ValidatorStatisticsTypeLinkedNodes, zero)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE validators
						SET
							linked_nodes =  (
							 	SELECT COALESCE((SELECT amount
								 FROM validator_statistics
								 WHERE validator_id = $1 AND statistic_type = $2 
								 ORDER BY block_height DESC LIMIT 1 ), 0))
						WHERE validator_id = $1`,
		params.ValidatorID, structs.ValidatorStatisticsTypeLinkedNodes)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

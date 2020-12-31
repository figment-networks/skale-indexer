package postgresql

import (
	"context"
	"database/sql"
	"math/big"
	"strconv"
	"strings"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

func (d *Driver) SaveValidatorStatistic(ctx context.Context, validatorID *big.Int, blockHeight uint64, statisticsType structs.StatisticTypeVS, amount *big.Int) (err error) {
	// (lukanus): Update value in validator_statistics unless the value already exists
	_, err = d.db.ExecContext(ctx, `
	INSERT INTO validator_statistics (validator_id, block_height, statistics_type, amount)
		( SELECT $1, $2, $3, $4	WHERE
			NOT EXISTS (
					SELECT 1 FROM validator_statistics
					WHERE validator_id = $1 AND statistics_type = $3 AND block_height < $2 AND amount = $4
					ORDER BY block_height DESC LIMIT 1
					)
		)
		ON CONFLICT (validator_id, block_height, statistics_type)
		DO UPDATE SET amount = EXCLUDED.amount;
		`, validatorID.String(), blockHeight, statisticsType, amount.String())
	return nil
}

func (d *Driver) GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	q := `SELECT
			DISTINCT ON (validator_id, statistics_type) id, created_at, validator_id, amount, block_height, statistics_type
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

	var (
		vldId  uint64
		amount string
	)

	for rows.Next() {
		vs := structs.ValidatorStatistics{}
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

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx,
		`SELECT id, created_at, validator_id, amount, block_height, statistics_type
			FROM validator_statistics
			WHERE
				validator_id = $1 AND statistics_type = $2
			ORDER BY block_height DESC`, params.ValidatorId, params.StatisticsTypeVS)
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
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
										SELECT t1.validator_id, sum(t1.amount) AS amount, $1 AS block_height, $2 AS statistics_type
									FROM
											( SELECT DISTINCT ON (delegation_id) validator_id, delegation_id, block_height, state, amount
											FROM delegations
												WHERE
													validator_id = $3 AND
												block_height <=$4
											ORDER BY delegation_id, block_height DESC) t1
										WHERE  t1.state IN ($5, $6) GROUP BY t1.validator_id
									ON CONFLICT (validator_id, block_height, statistics_type)
									DO UPDATE SET amount = EXCLUDED.amount`,
		params.BlockHeight, structs.ValidatorStatisticsTypeTotalStake, params.ValidatorId, params.BlockHeight, structs.DelegationStateACCEPTED, structs.DelegationStateUNDELEGATION_REQUESTED)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE validators
						SET staked = (
							 	SELECT amount
								 FROM validator_statistics
								 WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 )
						WHERE validator_id = $4`,
		params.ValidatorId, structs.ValidatorStatisticsTypeTotalStake, params.BlockHeight, params.ValidatorId)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()

}

func (d *Driver) CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	// BUG(l): This is wrong, would give random results. It either has to be calculated from  state, or nodes

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO validator_statistics (validator_id, block_height, statistics_type, amount)
				(SELECT validator_id, $1 AS block_height, $2 AS statistics_type, count(*) AS amount,
				FROM nodes
				WHERE validator_id = $3 AND status = $4
				GROUP BY validator_id LIMIT 1)
			ON CONFLICT (validator_id, block_height, statistics_type)
			DO UPDATE SET amount = EXCLUDED.amount `,
		params.BlockHeight, structs.ValidatorStatisticsTypeActiveNodes, params.ValidatorId, structs.NodeStatusActive)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE validators
						SET
							active_nodes = (SELECT amount
									FROM validator_statistics
									WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3)
						WHERE validator_id = $1`,
		params.ValidatorId, structs.ValidatorStatisticsTypeActiveNodes, params.BlockHeight)
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
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount, $1 AS block_height, $2 AS statistics_type
									FROM nodes
									WHERE validator_id = $3
									GROUP BY validator_id LIMIT 1)
								ON CONFLICT (validator_id, block_height, statistics_type)
								DO UPDATE SET amount = EXCLUDED.amount`,
		params.BlockHeight, structs.ValidatorStatisticsTypeLinkedNodes, params.ValidatorId, structs.NodeStatusActive)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(`UPDATE validators
						SET
							linked_nodes = (SELECT amount
											FROM validator_statistics
											WHERE validator_id = $1 AND statistics_type = $2 AND block_height = $3 )
						WHERE validator_id = $4`,
		params.ValidatorId, structs.ValidatorStatisticsTypeLinkedNodes, params.BlockHeight)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	getByStatementForVS = `SELECT d.id, d.created_at, d.updated_at, d.validator_id, d.amount, d.statistics_type FROM validator_statistics d WHERE d.statistics_type = $1 `
	byIdForVS           = `AND d.id = $2 `
	byValidatorIdForVS  = `AND d.validator_id = $2 `
	byStatusForVS       = `AND d.status = $3 `
	orderByCreatedAtVS  = `ORDER BY d.created_at DESC `
	calculateTotalStake = `INSERT INTO validator_statistics (updated_at, validator_id, amount, statistics_type) 
									(SELECT NOW(), validator_id, sum(amount) AS amount, $1 AS statistics_type FROM delegations
									WHERE validator_id = $2 AND status IN ($3 ,$4) GROUP BY validator_id)`
	calculateActiveNodes = `INSERT INTO validator_statistics (updated_at, validator_id, amount, statistics_type) 
									(SELECT NOW(), validator_id, count(*) AS amount, $1 AS statistics_type FROM nodes
									WHERE validator_id = $2 AND status = $3 GROUP BY validator_id)`
	calculateLinkedNodes = `INSERT INTO validator_statistics (updated_at, validator_id, amount, statistics_type) 
									(SELECT NOW(), validator_id, count(*) AS amount, $1 AS statistics_type FROM nodes
									WHERE validator_id = $2 GROUP BY validator_id)`
	updateTotalStake  = `UPDATE validators SET updated_at = NOW(), staked = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 ORDER BY created_at DESC LIMIT 1) WHERE validator_id = $3`
	updateActiveNodes = `UPDATE validators SET updated_at = NOW(), active_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 ORDER BY created_at DESC LIMIT 1) WHERE validator_id = $3`
	updateLinkedNodes = `UPDATE validators SET updated_at = NOW(), linked_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 ORDER BY created_at DESC LIMIT 1) WHERE validator_id = $3`
)

func (d *Driver) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForVS, byIdForVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS, params.Id)
	} else if params.ValidatorId > 0 && params.Status == 0 {
		q = fmt.Sprintf("%s%s%s", getByStatementForVS, byValidatorIdForVS, orderByCreatedAtVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS, params.ValidatorId)
	} else if params.ValidatorId > 0 && params.Status > 0 {
		q = fmt.Sprintf("%s%s%s%s", getByStatementForVS, byValidatorIdForVS, byStatusForVS, orderByCreatedAtVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS, params.ValidatorId, params.Status)
	} else if params.StatisticTypeVS != 0 && params.Id == "" && params.ValidatorId == 0 && params.Status == 0 {
		q = fmt.Sprintf("%s%s", getByStatementForVS, orderByCreatedAtVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS)
	} else {
		return validatorStatistics, handler.ErrMissingParameter
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		d := structs.ValidatorStatistics{}
		err = rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ValidatorId, &d.Amount, &d.StatisticType)
		if err != nil {
			return nil, err
		}
		validatorStatistics = append(validatorStatistics, d)
	}
	if len(validatorStatistics) == 0 {
		return nil, handler.ErrNotFound
	}
	return validatorStatistics, nil
}

// update "TOTAL STAKE" = Accepted + UndelegationRequested
func (d *Driver) CalculateTotalStake(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateTotalStake, structs.TotalStakeStatisticsTypeVS, params.ValidatorId, structs.Accepted, structs.UndelegationRequested)
	if err == nil {
		_, err = d.db.Exec(updateTotalStake, params.ValidatorId, structs.TotalStakeStatisticsTypeVS, params.ValidatorId)
	}
	return err
}

// update "ACTIVE NODES"
func (d *Driver) CalculateActiveNodes(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateActiveNodes, structs.ActiveNodesStatisticsTypeVS, params.ValidatorId, structs.Active)
	if err == nil {
		_, err = d.db.Exec(updateActiveNodes, params.ValidatorId, structs.ActiveNodesStatisticsTypeVS, params.ValidatorId)
	}
	return err
}

// update "LINKED NODES"
func (d *Driver) CalculateLinkedNodes(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateLinkedNodes, structs.LinkedNodesStatisticsTypeVS, params.ValidatorId)
	if err == nil {
		_, err = d.db.Exec(updateLinkedNodes, params.ValidatorId, structs.LinkedNodesStatisticsTypeVS, params.ValidatorId)
	}
	return err
}

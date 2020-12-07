package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	getByStatementForDS        = `SELECT d.id, d.created_at, d.updated_at, d.validator_id, d.status, d.amount, d.statistics_type FROM delegation_statistics d WHERE d.statistics_type = $1 `
	byIdForDS                  = `AND d.id = $2 `
	byValidatorIdForDS         = `AND d.validator_id = $2 `
	byStatusForDS              = `AND d.status = $3 `
	orderByCreatedAtDS         = `ORDER BY d.created_at DESC `
	calculateLatestStatesForDS = `INSERT INTO delegation_statistics (updated_at, validator_id, status, amount, statistics_type) 
									(SELECT NOW(), validator_id, status, sum(amount) AS amount, $1 AS statistics_type FROM delegations
									WHERE validator_id = $2 GROUP BY validator_id, status)`
	getLatestDelegationStatesByValidatorForDS = `SELECT DISTINCT ON (statistics_type, validator_id, status) id, created_at, updated_at, validator_id, status, amount, statistics_type FROM delegation_statistics 
										WHERE statistics_type = $1 and validator_id = $2 ORDER BY statistics_type, validator_id, status, created_at DESC`
)

func (d *Driver) GetDelegationStatistics(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForDS, byIdForDS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeDS, params.Id)
	} else if params.ValidatorId > 0 && params.Status == 0 {
		q = fmt.Sprintf("%s%s%s", getByStatementForDS, byValidatorIdForDS, orderByCreatedAtDS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeDS, params.ValidatorId)
	} else if params.ValidatorId > 0 && params.Status > 0 {
		q = fmt.Sprintf("%s%s%s%s", getByStatementForDS, byValidatorIdForDS, byStatusForDS, orderByCreatedAtDS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeDS, params.ValidatorId, params.Status)
	} else if params.StatisticTypeDS != 0 && params.Id == "" && params.ValidatorId == 0 && params.Status == 0 {
		q = fmt.Sprintf("%s%s", getByStatementForDS, orderByCreatedAtDS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeDS)
	} else {
		return delegationStatistics, handler.ErrMissingParameter
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		d := structs.DelegationStatistics{}
		err = rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ValidatorId, &d.Status, &d.Amount, &d.StatisticType)
		if err != nil {
			return nil, err
		}
		delegationStatistics = append(delegationStatistics, d)
	}
	if len(delegationStatistics) == 0 {
		return nil, handler.ErrNotFound
	}
	return delegationStatistics, nil
}

// update "DELEGATION STATES"
func (d *Driver) CalculateLatestDelegationStatesStatistics(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateLatestStatesForDS, structs.StatesStatisticsTypeDS, params.ValidatorId)
	return err
}

// get "DELEGATION STATES"
func (d *Driver) GetLatestDelegationStates(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error) {
	rows, err := d.db.QueryContext(ctx, getLatestDelegationStatesByValidatorForDS, structs.StatesStatisticsTypeDS, params.ValidatorId)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		d := structs.DelegationStatistics{}
		err = rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ValidatorId, &d.Status, &d.Amount, &d.StatisticType)
		if err != nil {
			return nil, err
		}
		delegationStatistics = append(delegationStatistics, d)
	}
	if len(delegationStatistics) == 0 {
		return nil, handler.ErrNotFound
	}
	return delegationStatistics, nil
}

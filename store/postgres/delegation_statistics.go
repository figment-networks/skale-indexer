package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	getByStatementForDS = `SELECT d.id, d.created_at, d.updated_at, d.validator_id, d.status, d.amount, d.statistics_type FROM delegation_statistics d WHERE d.statistics_type = $1 `
	byIdForDS           = `AND d.id = $2 `
	byValidatorIdForDS  = `AND d.validator_id = $2 `
)

func (d *Driver) GetDelegationStatistics(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForDS, byIdForDS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticType, params.Id)
	} else if params.ValidatorId > 0 {
		q = fmt.Sprintf("%s%s", getByStatementForDS, byValidatorIdForDS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticType, params.ValidatorId)
	} else {
		rows, err = d.db.QueryContext(ctx, getByStatementForDS, params.StatisticType)
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


/*
update delegation states

INSERT INTO delegation_statistics  (updated_at, validator_id, status, amount, statistics_type)
(SELECT NOW(), 2 as validator_id, 1 as  status,  24 AS amount, 3 AS statistics_type  FROM delegations
WHERE validator_id = 2 GROUP BY)


fetch latest delegation states
SELECT DISTINCT ON (statistics_type, validator_id, status) *
FROM delegation_statistics
ORDER BY statistics_type, validator_id, status, created_at DESC ;

*/


/*
	list proposed ones
 */
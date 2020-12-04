package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	getByStatementForVS = `SELECT d.id, d.created_at, d.updated_at, d.validator_id, d.status, d.amount, d.statistics_type FROM validator_statistics d WHERE d.statistics_type = $1 `
	byIdForVS           = `AND d.id = $2 `
	byValidatorIdForVS  = `AND d.validator_id = $2 `
	byStatusForVS       = `AND d.status = $3 `
	orderByCreatedAtVS  = `ORDER BY d.created_at DESC `
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
		err = rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ValidatorId, &d.Status, &d.Amount, &d.StatisticType)
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

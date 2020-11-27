package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/structs"
)

const (
	getByStatementForDSS = `SELECT d.id, d.created_at, d.updated_at, d.validator_id, d.status, d.amount FROM delegation_state_statistics d `
	byIdForDSS           = `WHERE d.id =  $1 `
	byValidatorIdForDSS  = `WHERE d.validator_id =  $1 `
)

func (d *Driver) GetDelegationStateStatistics(ctx context.Context, params structs.QueryParams) (delegationStateStatistics []structs.DelegationStateStatistics, err error) {
	var q string
	var rows *sql.Rows
	if params.Id != "" {
		q = fmt.Sprintf("%s%s", getByStatementForDSS, byIdForDSS)
		rows, err = d.db.QueryContext(ctx, q, params.Id)
	} else if params.ValidatorId > 0 {
		q = fmt.Sprintf("%s%s", getByStatementForDSS, byValidatorIdForDSS)
		rows, err = d.db.QueryContext(ctx, q, params.ValidatorId)
	} else {
		rows, err = d.db.QueryContext(ctx, getByStatementForDSS)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		d := structs.DelegationStateStatistics{}
		err = rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, d.ValidatorId, d.Status, d.Amount)
		if err != nil {
			return nil, err
		}
		delegationStateStatistics = append(delegationStateStatistics, d)
	}
	if len(delegationStateStatistics) == 0 {
		return nil, handler.ErrNotFound
	}
	return delegationStateStatistics, nil
}

package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// TODO: run explain analyze to check full scan and add required indexes
const (
	getByStatementVS    = `SELECT d.id, d.created_at, d.updated_at, d.validator_id, d.amount, d.block_height, d.statistics_type FROM validator_statistics d WHERE d.statistics_type = $1 `
	byIdVS              = `AND d.id = $2 `
	byValidatorIdVS     = `AND d.validator_id = $2 `
	orderByCreatedAtVS  = `ORDER BY d.created_at DESC `
	calculateTotalStake = `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
 								SELECT t1.validator_id, sum(t1.amount) AS amount, $1 AS block_height, $2 AS statistics_type FROM
 									(SELECT  DISTINCT ON (delegation_id) validator_id, delegation_id, block_height, state, amount FROM delegations
 											WHERE validator_id = $3 AND block_height <=$4 ORDER BY delegation_id, block_height DESC)  t1
 								WHERE  t1.state IN ($5, $6) GROUP BY t1.validator_id`
	calculateActiveNodes = `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount, $1 AS block_height, $2 AS statistics_type FROM nodes
									WHERE validator_id = $3 AND state = $4 GROUP BY validator_id, delegation_id ORDER BY block_height DESC)`
	calculateLinkedNodes = `INSERT INTO validator_statistics (validator_id, amount, block_height, statistics_type)
									(SELECT validator_id, count(*) AS amount, $1 AS block_height, $2 AS statistics_type FROM nodes
									WHERE validator_id = $3 GROUP BY validator_id, delegation_id ORDER BY block_height DESC)`
	updateTotalStake  = `UPDATE validators SET updated_at = NOW(), staked = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height <=$3 ORDER BY block_height DESC LIMIT 1) WHERE validator_id = $3`
	updateActiveNodes = `UPDATE validators SET updated_at = NOW(), active_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height <=$3 ORDER BY block_height DESC LIMIT 1) WHERE validator_id = $3`
	updateLinkedNodes = `UPDATE validators SET updated_at = NOW(), linked_nodes = (SELECT amount FROM validator_statistics WHERE validator_id = $1 AND statistics_type = $2 AND block_height <=$3 ORDER BY block_height DESC LIMIT 1) WHERE validator_id = $3`
)

func (d *Driver) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	var q string
	var rows *sql.Rows
	if params.ValidatorId != 0 {
		q = fmt.Sprintf("%s%s%s", getByStatementVS, byValidatorIdVS, orderByCreatedAtVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS, params.ValidatorId)
	} else if params.Id == "" && params.ValidatorId == 0 {
		q = fmt.Sprintf("%s%s", getByStatementVS, orderByCreatedAtVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS)
	} else {
		q = fmt.Sprintf("%s%s", getByStatementVS, byIdVS)
		rows, err = d.db.QueryContext(ctx, q, params.StatisticTypeVS, params.Id)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		d := structs.ValidatorStatistics{}
		err = rows.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.ValidatorId, &d.Amount, &d.BlockHeight, &d.StatisticType)
		if err != nil {
			return nil, err
		}
		validatorStatistics = append(validatorStatistics, d)
	}
	return validatorStatistics, nil
}

// update "TOTAL STAKE" = Accepted + UndelegationRequested
func (d *Driver) CalculateTotalStake(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateTotalStake, params.ETHBlockHeight, structs.ValidatorStatisticsTypeTotalStake, params.ValidatorId, params.ETHBlockHeight, structs.DelegationStateACCEPTED, structs.DelegationStateUNDELEGATION_REQUESTED)
	if err == nil {
		_, err = d.db.Exec(updateTotalStake, params.ValidatorId, structs.ValidatorStatisticsTypeTotalStake, params.ValidatorId, params.ETHBlockHeight)
	}
	return err
}

// update "ACTIVE NODES"
// TODO: change calculation active nodes = active + leaving nodes
func (d *Driver) CalculateActiveNodes(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateActiveNodes, params.ETHBlockHeight, structs.ValidatorStatisticsTypeActiveNodes, params.ValidatorId, structs.NodeStatusActive)
	if err == nil {
		_, err = d.db.Exec(updateActiveNodes, params.ValidatorId, structs.ValidatorStatisticsTypeActiveNodes, params.ValidatorId, params.ETHBlockHeight)
	}
	return err
}

// update "LINKED NODES"
func (d *Driver) CalculateLinkedNodes(ctx context.Context, params structs.QueryParams) error {
	_, err := d.db.Exec(calculateLinkedNodes, params.ETHBlockHeight, structs.ValidatorStatisticsTypeLinkedNodes, params.ValidatorId)
	if err == nil {
		_, err = d.db.Exec(updateLinkedNodes, params.ValidatorId, structs.ValidatorStatisticsTypeLinkedNodes, params.ValidatorId, params.ETHBlockHeight)
	}
	return err
}

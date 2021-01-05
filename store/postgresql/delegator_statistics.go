package postgresql

import (
	"context"
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"math/big"
	"strconv"
	"strings"
)

// SaveDelegatorStatistics saves delegator statistics
func (d *Driver) SaveDelegatorStatistics(ctx context.Context, vs structs.DelegatorStatistics) error {
	_, err := d.db.Exec(`INSERT INTO delegator_statistics ("holder", "amount", "block_height", "statistics_type") 
			VALUES ($1, $2, $3, $4) 
			ON CONFLICT (statistics_type, holder, block_height)
			DO UPDATE SET
			amount = EXCLUDED.amount `,
		vs.Holder.Hash().Big().String(),
		vs.Amount.String(),
		vs.BlockHeight,
		vs.StatisticsTypeDS)
	return err
}

func (d *Driver) GetDelegatorStatistics(ctx context.Context, params structs.DelegatorStatisticsParams) (delegatorStatistics []structs.DelegatorStatistics, err error) {
	q := `SELECT
			DISTINCT ON (holder, statistics_type) 
				id, created_at, holder, amount, block_height, statistics_type 
			FROM delegator_statistics `
	var (
		args   []interface{}
		wherec []string
		i      = 1
	)
	if params.Holder != "" {
		wherec = append(wherec, ` holder =  $`+strconv.Itoa(i))
		args = append(args, common.HexToAddress(params.Holder).Hash().Big().String())
		i++
	}
	if params.StatisticsTypeDS != "" {
		wherec = append(wherec, ` statistics_type =  $`+strconv.Itoa(i))
		args = append(args, params.StatisticsTypeDS)
		i++
	}
	if len(args) > 0 {
		q += ` WHERE `
	}
	q += strings.Join(wherec, " AND ")
	q += ` ORDER BY holder ASC, statistics_type, block_height DESC`

	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vs := structs.DelegatorStatistics{}
		var holder []byte
		var amount string
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &holder, &amount, &vs.BlockHeight, &vs.StatisticsTypeDS)
		if err != nil {
			return nil, err
		}
		h := new(big.Int)
		h.SetString(string(holder), 10)
		vs.Holder.SetBytes(h.Bytes())
		amnt := new(big.Int)
		amnt.SetString(amount, 10)
		vs.Amount = amnt
		delegatorStatistics = append(delegatorStatistics, vs)
	}
	return delegatorStatistics, nil
}

func (d *Driver) GetDelegatorStatisticsChart(ctx context.Context, params structs.DelegatorStatisticsParams) (delegatorStatistics []structs.DelegatorStatistics, err error) {
	q := `SELECT id, created_at, holder, amount, block_height, statistics_type 
			FROM delegator_statistics WHERE holder = $1 AND statistics_type = $2 ORDER BY block_height DESC`
	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, common.HexToAddress(params.Holder).Hash().Big().String(), params.StatisticsTypeDS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vs := structs.DelegatorStatistics{}
		var holder []byte
		var amount string
		err = rows.Scan(&vs.ID, &vs.CreatedAt, &holder, &amount, &vs.BlockHeight, &vs.StatisticsTypeDS)
		if err != nil {
			return nil, err
		}
		h := new(big.Int)
		h.SetString(string(holder), 10)
		vs.Holder.SetBytes(h.Bytes())
		amnt := new(big.Int)
		amnt.SetString(amount, 10)
		vs.Amount = amnt
		delegatorStatistics = append(delegatorStatistics, vs)
	}
	return delegatorStatistics, nil
}

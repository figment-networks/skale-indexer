package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
)

// SaveSystemEvent saves system events
func (d *Driver) SaveSystemEvent(ctx context.Context, se structs.SystemEvent) error {

	_, err := d.db.Exec(
		`INSERT INTO system_events(
			"height",
			"kind",
			"time",
			"sender",
			"recipient",
			"sender_id",
			"recipient_id",
			"before",
			"after",
			"change"
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT ( height, kind, sender, sender_id, recipient, recipient_id)
		DO UPDATE SET
			before = EXCLUDED.before,
			after = EXCLUDED.after,
			change = EXCLUDED.change
		`,
		se.Height,
		se.Kind,
		se.Time,
		se.Sender.Hash().Big().String(),
		se.Recipient.Hash().Big().String(),
		se.SenderID.String(),
		se.RecipientID.String(),
		se.Before.String(),
		se.After.String(),
		se.Change.String(),
	)
	return err
}

// GetSystemEvents gets contract events
func (d *Driver) GetSystemEvents(ctx context.Context, params structs.SystemEventParams) (systemEvents []structs.SystemEvent, err error) {

	q := `SELECT height, kind, time, sender, sender_id, recipient, recipient_id, before, after, change FROM system_events `

	var (
		args   []interface{}
		whereC []string
		i      = 1
	)

	if params.Address != "" {
		whereC = append(whereC, `( sender = $`+strconv.Itoa(i)+` OR recipient =  $`+strconv.Itoa(i)+` )`)
		args = append(args, common.HexToAddress(params.Address).Hash().Big().String())
		i++
	}

	if params.ReceiverID > 0 {
		whereC = append(whereC, ` recipient_id = $`+strconv.Itoa(i))
		args = append(args, params.ReceiverID)
		i++
	}

	if params.SenderID > 0 {
		whereC = append(whereC, ` sender_id = $`+strconv.Itoa(i))
		args = append(args, params.SenderID)
		i++
	}

	if params.Kind != "" {
		whereC = append(whereC, ` kind = $`+strconv.Itoa(i))
		args = append(args, params.Kind)
		i++
	}

	if params.After > 0 {
		whereC = append(whereC, ` height > $`+strconv.Itoa(i))
		args = append(args, params.After)
		i++
	}

	if len(whereC) > 0 {
		q += " WHERE "
	}

	q += strings.Join(whereC, " AND ")
	q += `ORDER BY height DESC `

	if params.Limit > 0 {
		q += " LIMIT " + strconv.FormatUint(uint64(params.Limit), 10)
		if params.Offset > 0 {
			q += " OFFSET " + strconv.FormatUint(uint64(params.Offset), 10)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	var rows *sql.Rows
	rows, err = d.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var (
		sender    string
		recipient string

		senderID    string
		recipientID string

		beforeValue string
		afterValue  string

		change string
	)
	for rows.Next() {
		e := structs.SystemEvent{}
		if err = rows.Scan(&e.Height, &e.Kind, &e.Time, &sender, &senderID, &recipient, &recipientID, &beforeValue, &afterValue, &change); err != nil {
			return nil, err
		}

		p := new(big.Int)
		p.SetString(string(beforeValue), 10)
		e.Before.SetBytes(p.Bytes())

		p.SetString(string(afterValue), 10)
		e.After.SetBytes(p.Bytes())

		p.SetString(string(sender), 10)
		e.Sender.SetBytes(p.Bytes())

		p.SetString(string(sender), 10)
		e.Sender.SetBytes(p.Bytes())

		p.SetString(string(recipient), 10)
		e.Recipient.SetBytes(p.Bytes())

		p.SetString(string(senderID), 10)
		e.SenderID.SetBytes(p.Bytes())

		p.SetString(string(recipientID), 10)
		e.RecipientID.SetBytes(p.Bytes())

		f := new(big.Float)
		f.SetString(change)
		e.Change = *f

		systemEvents = append(systemEvents, e)
	}

	return systemEvents, nil
}

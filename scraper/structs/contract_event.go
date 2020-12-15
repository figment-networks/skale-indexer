package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"time"
)

type ContractEvent struct {
	ID              string         `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       pq.NullTime    `json:"updated_at"`
	ContractName    string         `json:"contract_name"`
	EventName       string         `json:"event_name"`
	ContractAddress common.Address `json:"contract_address"`
	BlockHeight     uint64         `json:"block_height"`
	Time            time.Time      `json:"time"`
	TransactionHash common.Hash    `json:"transaction_hash"`
	Removed         bool           `json:"removed"`
	Params          PropertyMap    `json:"params"`
	BoundType       string         `json:"bound_type"`
	BoundId         []BoundId      `json:"bound_id"`
	BoundAddress    common.Address `json:"bound_address"`
}

type BoundId uint64

func (a *BoundId) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type PropertyMap map[string]interface{}

func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type ContractEvent struct {
	ID              string         `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	ContractName    string         `json:"contract_name"`
	EventName       string         `json:"event_name"`
	ContractAddress common.Address `json:"contract_address"`
	BlockHeight     uint64         `json:"block_height"`
	Time            time.Time      `json:"time"`
	TransactionHash common.Hash    `json:"transaction_hash"`
	Removed         bool           `json:"removed"`
	Params          PropertyMap    `json:"params"`
	BoundType       string         `json:"bound_type"`
	BoundId         big.Int      `json:"bound_id"`
	BoundAddress    common.Address `json:"bound_address"`
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
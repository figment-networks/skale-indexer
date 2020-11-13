package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

type Time struct {
	time.Time
}

func NewTimeFromTime(t time.Time) *Time {
	return &Time{
		Time: t,
	}
}

func NewTimeFromTimestamp(timestamp timestamp.Timestamp) *Time {
	return &Time{
		Time: time.Unix(timestamp.GetSeconds(), int64(timestamp.GetNanos())),
	}
}

func (t *Time) Duration(m Time) int64 {
	return t.Sub(m.Time).Milliseconds()
}

func (t *Time) Equal(m Time) bool {
	return t.Time.Equal(m.Time)
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *Time) Scan(value interface{}) error {
	tm, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	t.Time = tm
	return nil
}

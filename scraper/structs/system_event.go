package structs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type SystemEvent struct {
	ID          string         `json:"id"`
	Height      uint64         `json:"height"`
	Time        time.Time      `json:"time"`
	Kind        SysEvtType     `json:"kind"`
	SenderID    big.Int        `json:"sender_id"`
	RecipientID big.Int        `json:"recipient_id"`
	Sender      common.Address `json:"sender"`
	Recipient   common.Address `json:"recipient"`

	Before big.Int   `json:"before"`
	After  big.Int   `json:"after"`
	Change big.Float `json:"change"`
}

type SysEvtType uint64

const (
	SysEvtTypeNewDelegation SysEvtType = iota + 1
	SysEvtTypeDelegationAccepted
	SysEvtTypeDelegationRejected
	SysEvtTypeUndeledationRequested
	SysEvtTypeJoinedActiveSet
	SysEvtTypeLeftActiveSet
	SysEvtTypeSlashed
	SysEvtTypeForgiven
	SysEvtTypeMDRChanged
	SysEvtTypeFeeChanged
)

var (
	SysEvtTypes = map[SysEvtType]string{
		SysEvtTypeNewDelegation:         "new_delegation",
		SysEvtTypeDelegationAccepted:    "delegation_accepted",
		SysEvtTypeDelegationRejected:    "delegation_rejected",
		SysEvtTypeUndeledationRequested: "undelegation_requested",
		SysEvtTypeJoinedActiveSet:       "joined_active_set",
		SysEvtTypeLeftActiveSet:         "left_active_set",
		SysEvtTypeMDRChanged:            "mdr_change",
		SysEvtTypeFeeChanged:            "fee_change",
		SysEvtTypeSlashed:               "slashed",
		SysEvtTypeForgiven:              "forgiven",
	}
)

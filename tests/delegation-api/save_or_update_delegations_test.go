package tests

import (
	"bytes"
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	invalidSyntaxForDelegations = `[{
        "holder_invalid": 1,
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": "2014-11-12T11:45:26.371Z",
        "started": "2014-11-12T11:45:26.371Z",
        "finished": "2014-11-12T11:45:26.371Z",
        "info": "info1"
	`
	invalidPropertyNameForDelegations = `[{
        "holder_invalid": 1,
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": "2014-11-12T11:45:26.371Z",
        "started": "2014-11-12T11:45:26.371Z",
        "finished": "2014-11-12T11:45:26.371Z",
        "info": "info1"
    }]`
	validJsonForDelegations = `[{
        "holder": 1,
        "validator_id": 2,
		"skale_id": 11,
		"eth_block_height": 1000,
        "amount": 0,
        "delegation_period": 0,
        "created": "2014-11-12T11:45:26.371Z",
        "started": "2014-11-12T11:45:26.371Z",
        "finished": "2014-11-12T11:45:26.371Z",
        "info": "info1",
        "status": 1,
        "smart_contract_index": 1903,
        "smart_contract_address": 1001
    },	
	{
        "holder": 2,
        "validator_id": 2,
		"skale_id": 12,
		"eth_block_height": 1001,
        "amount": 0,
        "delegation_period": 0,
        "created": "2014-11-12T11:45:26.371Z",
        "started": "2014-11-12T11:45:26.371Z",
        "finished": "2014-11-12T11:45:26.371Z",
        "info": "info2",
        "status": 2,
        "smart_contract_index": 1904,
        "smart_contract_address": 1002
    }	
	]`
	// same value should be used in json examples above for valid cases
	DummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleDelegations []structs.Delegation

func TestSaveOrUpdateDelegations(t *testing.T) {
	var holder uint64 = 1
	var validatorId uint64 = 2
	var amount uint64 = 0
	var delegationPeriod uint64 = 0
	exampleTime, _ := time.Parse(handler.Layout, DummyTime)
	var created = exampleTime
	var started = exampleTime
	var finished = exampleTime
	info := "info1"
	exampleDelegation := structs.Delegation{
		Holder:               holder,
		ValidatorId:          validatorId,
		SkaleId:              uint64(11),
		ETHBlockHeight:       uint64(1000),
		Amount:               amount,
		DelegationPeriod:     delegationPeriod,
		Created:              created,
		Started:              started,
		Finished:             finished,
		Info:                 info,
		Status:               1,
		SmartContractIndex:   1903,
		SmartContractAddress: 1001,
	}
	var holder2 uint64 = 2
	var validatorId2 uint64 = 2
	var amount2 uint64 = 0
	var delegationPeriod2 uint64 = 0
	var created2 = exampleTime
	var started2 = exampleTime
	var finished2 = exampleTime
	info2 := "info2"
	exampleDelegation2 := structs.Delegation{
		Holder:               holder2,
		ValidatorId:          validatorId2,
		SkaleId:              uint64(12),
		ETHBlockHeight:       uint64(1001),
		Amount:               amount2,
		DelegationPeriod:     delegationPeriod2,
		Created:              created2,
		Started:              started2,
		Finished:             finished2,
		Info:                 info2,
		Status:               2,
		SmartContractIndex:   1904,
		SmartContractAddress: 1002,
	}
	exampleDelegations = append(exampleDelegations, exampleDelegation)
	exampleDelegations = append(exampleDelegations, exampleDelegation2)

	tests := []struct {
		number      int
		name        string
		req         *http.Request
		delegations []structs.Delegation
		dbResponse  error
		code        int
	}{
		{
			number: 1,
			name:   "not allowed method",
			req: &http.Request{
				Method: http.MethodGet,
			},
			code: http.StatusMethodNotAllowed,
		},
		{
			number: 2,
			name:   "invalid json syntax request body",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForDelegations))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "missing parameter",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForDelegations))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForDelegations))),
			},
			delegations: exampleDelegations,
			dbResponse:  errors.New("internal error"),
			code:        http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForDelegations))),
			},
			delegations: exampleDelegations,
			code:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().SaveOrUpdateDelegations(tt.req.Context(), tt.delegations).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateDelegations)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

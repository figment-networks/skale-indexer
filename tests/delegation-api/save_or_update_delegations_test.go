package tests

import (
	"../../client"
	"../../handler"
	"../../store"
	"../../structs"
	"bytes"
	"errors"
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
        "holder_invalid": "holder1",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": 0
        "started": 0,
        "finished": 0,
        "info": "info1"
	`
	invalidPropertyNameForDelegations = `[{
        "holder_invalid": "holder1",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": 0,
        "started": 0,
        "finished": 0,
        "info": "info1"
    }]`
	validJsonForDelegations = `[{
        "holder": "holder1",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": "2014-11-12T11:45:26.371Z",
        "started": "2014-11-12T11:45:26.371Z",
        "finished": "2014-11-12T11:45:26.371Z",
        "info": "info1"
    },	
	{
        "holder": "holder2",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": "2014-11-12T11:45:26.371Z",
        "started": "2014-11-12T11:45:26.371Z",
        "finished": "2014-11-12T11:45:26.371Z",
        "info": "info2"
    }	
	]`
	// same value should be used in json examples above for valid cases
	dummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleDelegations []structs.Delegation

func TestSaveOrUpdateDelegations(t *testing.T) {
	holder := "holder1"
	var validatorId uint64 = 2
	var amount uint64 = 0
	var delegationPeriod uint64 = 0
	layout := "2006-01-02T15:04:05.000Z"
	exampleTime, _ := time.Parse(layout, dummyTime)
	var created time.Time = exampleTime
	var started time.Time = exampleTime
	var finished time.Time = exampleTime
	info := "info1"
	exampleDelegation := structs.Delegation{
		Holder:           &holder,
		ValidatorId:      &validatorId,
		Amount:           &amount,
		DelegationPeriod: &delegationPeriod,
		Created:          &created,
		Started:          &started,
		Finished:         &finished,
		Info:             &info,
	}
	holder2 := "holder2"
	var validatorId2 uint64 = 2
	var amount2 uint64 = 0
	var delegationPeriod2 uint64 = 0
	var created2 time.Time = exampleTime
	var started2 time.Time = exampleTime
	var finished2 time.Time = exampleTime
	info2 := "info2"
	exampleDelegation2 := structs.Delegation{
		Holder:           &holder2,
		ValidatorId:      &validatorId2,
		Amount:           &amount2,
		DelegationPeriod: &delegationPeriod2,
		Created:          &created2,
		Started:          &started2,
		Finished:         &finished2,
		Info:             &info2,
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
			name:   "bad request",
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
			if tt.number == 4 || tt.number == 5 {
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

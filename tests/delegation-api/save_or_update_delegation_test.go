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
)

const (
	invalidSyntaxForDelegation = `{
        "holder_invalid": "holder1",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": 0,
        "started": 0,
        "finished": 0,
        "info": "info1",
	`
	invalidPropertyNameForDelegation = `{
        "holder_invalid": "holder1",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": 0,
        "started": 0,
        "finished": 0,
        "info": "info1"
    }`
	validJsonForDelegation = `{
        "holder": "holder1",
        "validator_id": 2,
        "amount": 0,
        "delegation_period": 0,
        "created": 0,
        "started": 0,
        "finished": 0,
        "info": "info1"
    }`
)

var exampleDelegation structs.Delegation

func TestSaveOrUpdateDelegation(t *testing.T) {
	holder := "holder1"
	var validatorId uint64 = 2
	var amount uint64 = 0
	var delegationPeriod uint64 = 0
	var created uint64 = 0
	var started uint64 = 0
	var finished uint64 = 0
	info := "info1"
	exampleDelegation = structs.Delegation{
		Holder:           &holder,
		ValidatorId:      &validatorId,
		Amount:           &amount,
		DelegationPeriod: &delegationPeriod,
		Created:          &created,
		Started:          &started,
		Finished:         &finished,
		Info:             &info,
	}
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		delegation structs.Delegation
		dbResponse error
		code       int
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
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForDelegation))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "bad request",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForDelegation))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForDelegation))),
			},
			delegation: exampleDelegation,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForDelegation))),
			},
			delegation: exampleDelegation,
			code:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number == 4 || tt.number == 5 {
				mockDB.EXPECT().SaveOrUpdateDelegation(tt.req.Context(), tt.delegation).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateDelegation)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

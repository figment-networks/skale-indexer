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
)

const (
	invalidSyntaxForValidators = `[{
        "name": "name_test",
        "address": "validator_address_test",
        "description": "description_test",
        "fee_rate": 1,
	`
	invalidPropertyNameForValidators = `[{
        "name_invalid": "name_test",
        "address": "validator_address_test",
        "description": "",
        "fee_rate": 1,
    }]`
	validJsonForValidators = `[{
        "name": "name_test",
        "address": [1,2],
        "description": "description_test",
        "fee_rate": 1,
        "active": true,
		"active_nodes": 2,
		"staked":  10,
		"pending": 15,
        "rewards": 20,
		"optional_info" : {
						"data": [	
							{
								"requested_address": "requested_address_test",
								"minimum_delegation_amount" : 2,
								"accept_new_requests": true,
								"trusted": true
							}	
						]	
			}
    },	
	{
     	"name": "name_test",
        "address": [1, 2],
        "description": "description_test",
        "fee_rate": 1,
        "active": false,
		"active_nodes": 0,
		"staked": 100,
		"pending": 0,
        "rewards": 140,
		"optional_info" : {
						"data": [	
							{
								"requested_address": "requested_address_test",
								"minimum_delegation_amount" : 2,
								"accept_new_requests": true,
								"trusted": true
							}	
						]	
			}
    }	
	]`
)

var exampleValidators []structs.Validator

func TestSaveOrUpdateDelegations(t *testing.T) {
	exampleValidator := structs.Validator{
		Name:        "name_test",
		Address:     []int{1, 2},
		Description: "description_test",
		FeeRate:     uint64(1),
		Active:      true,
		ActiveNodes: 2,
		Staked:      uint64(10),
		Pending:     uint64(15),
		Rewards:     uint64(20),
		OptionalInfo: structs.OptionalInfo{[]structs.Data{
			{
				RequestedAddress: "requested_address_test",
				MinimumDelegationAmount: 2,
				AcceptNewRequests: true,
				Trusted: true,
			},
		},},
	}
	exampleValidator2 := structs.Validator{
		Name:        "name_test",
		Address:     []int{1, 2},
		Description: "description_test",
		FeeRate:     uint64(1),
		Active:      false,
		ActiveNodes: 0,
		Staked:      uint64(100),
		Pending:     uint64(0),
		Rewards:     uint64(140),
		OptionalInfo: structs.OptionalInfo{ []structs.Data{
			{
				RequestedAddress: "requested_address_test",
				MinimumDelegationAmount: 2,
				AcceptNewRequests: true,
				Trusted: true,
			},
		},},
	}
	exampleValidators = append(exampleValidators, exampleValidator)
	exampleValidators = append(exampleValidators, exampleValidator2)

	tests := []struct {
		number     int
		name       string
		req        *http.Request
		validators []structs.Validator
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
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForValidators))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "missing parameter",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForValidators))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForValidators))),
			},
			validators: exampleValidators,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForValidators))),
			},
			validators: exampleValidators,
			code:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().SaveOrUpdateValidators(tt.req.Context(), tt.validators).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateValidators)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

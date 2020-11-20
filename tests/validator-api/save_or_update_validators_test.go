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
	invalidSyntaxForValidators = `[{
        "name": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": "2014-11-12T11:45:26.371Z",
        "minimum_delegation_amount": 0,
        "accept_new_requests": false
	`
	invalidPropertyNameForValidators = `[{
        "name_invalid": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": "2014-11-12T11:45:26.371Z",
        "minimum_delegation_amount": 0,
        "accept_new_requests": false
    }]`
	validJsonForValidators = `[{
       "name": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": "2014-11-12T11:45:26.371Z",
        "minimum_delegation_amount": 0,
        "accept_new_requests": false
    },	
	{
     	"name": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": "2014-11-12T11:45:26.371Z",
        "minimum_delegation_amount": 0,
        "accept_new_requests": false
    }	
	]`
	// same value should be used in json examples above for valid cases
	dummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleValidators []structs.Validator

func TestSaveOrUpdateDelegations(t *testing.T) {
	name := "name_test"
	validatorAddress := "validator_address_test"
	requestedAddress := "requested_address_test"
	description := "description_test"
	var feeRate uint64 = 1
	layout := "2006-01-02T15:04:05.000Z"
	exampleTime, _ := time.Parse(layout, dummyTime)
	var registrationTime time.Time = exampleTime
	var minimumDelegationAmount uint64 = 0
	var acceptNewRequests = false
	exampleValidator := structs.Validator{
		Name:                    name,
		ValidatorAddress:        validatorAddress,
		RequestedAddress:        requestedAddress,
		Description:             description,
		FeeRate:                 feeRate,
		RegistrationTime:        registrationTime,
		MinimumDelegationAmount: minimumDelegationAmount,
		AcceptNewRequests:       acceptNewRequests,
	}
	name2 := "name_test"
	validatorAddress2 := "validator_address_test"
	requestedAddress2 := "requested_address_test"
	description2 := "description_test"
	var feeRate2 uint64 = 1
	var registrationTime2 time.Time = exampleTime
	var minimumDelegationAmount2 uint64 = 0
	var acceptNewRequests2 = false
	exampleValidator2 := structs.Validator{
		Name:                    name2,
		ValidatorAddress:        validatorAddress2,
		RequestedAddress:        requestedAddress2,
		Description:             description2,
		FeeRate:                 feeRate2,
		RegistrationTime:        registrationTime2,
		MinimumDelegationAmount: minimumDelegationAmount2,
		AcceptNewRequests:       acceptNewRequests2,
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
			name:   "bad request",
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

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
	invalidSyntaxForValidator = `{
        "name": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": 0,
        "minimum_delegation_amount": 0,
        "accept_new_requests": false,
	`
	invalidPropertyNameForValidator = `{
        "name_invalid": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": 0,
        "minimum_delegation_amount": 0,
        "accept_new_requests": false
    }`
	validJsonForValidator = `{
        "name": "name_test",
        "validator_address": "validator_address_test",
        "requested_address": "requested_address_test",
        "description": "description_test",
        "fee_rate": 1,
        "registration_time": 0,
        "minimum_delegation_amount": 0,
        "accept_new_requests": false
    }`
)

var exampleValidator structs.Validator

func TestSaveOrUpdateValidator(t *testing.T) {
	name := "name_test"
	validatorAddress := "validator_address_test"
	requestedAddress := "requested_address_test"
	description := "description_test"
	var feeRate uint64 = 1
	var registrationTime uint64 = 0
	var minimumDelegationAmount uint64 = 0
	var acceptNewRequests = false
	exampleValidator = structs.Validator{
		Name:                    &name,
		ValidatorAddress:        &validatorAddress,
		RequestedAddress:        &requestedAddress,
		Description:             &description,
		FeeRate:                 &feeRate,
		RegistrationTime:        &registrationTime,
		MinimumDelegationAmount: &minimumDelegationAmount,
		AcceptNewRequests:       &acceptNewRequests,
	}

	tests := []struct {
		number     int
		name       string
		req        *http.Request
		validator  structs.Validator
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
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForValidator))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "bad request",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForValidator))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForValidator))),
			},
			validator:  exampleValidator,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForValidator))),
			},
			validator: exampleValidator,
			code:      http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number == 4 || tt.number == 5 {
				mockDB.EXPECT().SaveOrUpdateValidator(tt.req.Context(), tt.validator).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateValidator)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

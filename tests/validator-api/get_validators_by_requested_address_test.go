package tests

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var vldsByRequestedAddress = make([]structs.Validator, 1)

func TestGetValidatorsByRequestedAddress(t *testing.T) {
	name := "name_test"
	validatorAddress := "validator_address_test"
	requestedAddress := "requested_address_test"
	description := "description"
	var feeRate uint64 = 1
	layout := "2006-01-02T15:04:05.000Z"
	exampleTime, _ := time.Parse(layout, dummyTime)
	var registrationTime = exampleTime
	var minimumDelegationAmount uint64 = 0
	var acceptNewRequests = true
	vld := structs.Validator{
		Name:                    name,
		Address:                 validatorAddress,
		RequestedAddress:        requestedAddress,
		Description:             description,
		FeeRate:                 feeRate,
		RegistrationTime:        registrationTime,
		MinimumDelegationAmount: minimumDelegationAmount,
		AcceptNewRequests:       acceptNewRequests,
		Trusted:                 true,
	}
	vldsByRequestedAddress = append(vldsByRequestedAddress, vld)
	tests := []struct {
		number           int
		name             string
		req              *http.Request
		requestedAddress string
		validators       []structs.Validator
		dbResponse       error
		code             int
	}{
		{
			number: 1,
			name:   "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			code: http.StatusMethodNotAllowed,
		},
		{
			number: 2,
			name:   "missing parameter",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "empty parameter",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "requested_address=",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "requested_address=requested_address_test",
				},
			},
			requestedAddress: requestedAddress,
			dbResponse:       handler.ErrNotFound,
			code:             http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "requested_address=requested_address_test",
				},
			},
			requestedAddress: requestedAddress,
			dbResponse:       errors.New("internal error"),
			code:             http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "requested_address=requested_address_test",
				},
			},
			requestedAddress: requestedAddress,
			validators:       vldsByRequestedAddress,
			code:             http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetValidatorsByRequestedAddress(tt.req.Context(), tt.requestedAddress).Return(tt.validators, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorsByRequestedAddress)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

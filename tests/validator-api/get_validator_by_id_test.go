package tests

import (
	"../../client"
	"../../handler"
	"../../store"
	"../../structs"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var vldById structs.Validator

func TestGetValidatorById(t *testing.T) {
	name := "name_test"
	validatorAddress := "validator_address_test"
	requestedAddress := "requested_address_test"
	description := "description"
	var feeRate uint64 = 1
	layout := "2006-01-02T15:04:05.000Z"
	exampleTime, _ := time.Parse(layout, dummyTime)
	var registrationTime time.Time = exampleTime
	var minimumDelegationAmount uint64 = 0
	var acceptNewRequests = true
	vldById = structs.Validator{
		Name:                    name,
		ValidatorAddress:        validatorAddress,
		RequestedAddress:        requestedAddress,
		Description:             description,
		FeeRate:                 feeRate,
		RegistrationTime:        registrationTime,
		MinimumDelegationAmount: minimumDelegationAmount,
		AcceptNewRequests:       acceptNewRequests,
	}
	var id = "id_test"
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		id         string
		validator  structs.Validator
		dbResponse error
		code       int
	}{
		{
			number: 1,
			name:   "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			id:   id,
			code: http.StatusMethodNotAllowed,
		},
		{
			number: 2,
			name:   "bad request",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "invalid id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=",
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
					RawQuery: "id=id_test",
				},
			},
			id:         id,
			dbResponse: errors.New("record not found"),
			code:       http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=id_test",
				},
			},
			id:         id,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=id_test",
				},
			},
			id:        id,
			validator: vldById,
			code:      http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetValidatorById(tt.req.Context(), tt.id).Return(tt.validator, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorById)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

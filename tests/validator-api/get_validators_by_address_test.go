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
)

var vldsByValidatorAddress = make([]structs.Validator, 1)

func TestGetValidatorsByAddress(t *testing.T) {
	validatorAddress := "validator_address_test"
	vld := structs.Validator{
		Name:        "name_test",
		Address:     validatorAddress,
		Description: "description",
	}
	vldsByValidatorAddress = append(vldsByValidatorAddress, vld)
	tests := []struct {
		number           int
		name             string
		req              *http.Request
		validatorAddress string
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
					RawQuery: "address=",
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
					RawQuery: "address=validator_address_test",
				},
			},
			validatorAddress: validatorAddress,
			dbResponse:       handler.ErrNotFound,
			code:             http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "address=validator_address_test",
				},
			},
			validatorAddress: validatorAddress,
			dbResponse:       errors.New("internal error"),
			code:             http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "address=validator_address_test",
				},
			},
			validatorAddress: validatorAddress,
			validators:       vldsByValidatorAddress,
			code:             http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetValidatorsByAddress(tt.req.Context(), tt.validatorAddress).Return(tt.validators, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorsByAddress)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

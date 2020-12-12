package validator

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client/structs"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetValidatorsByValidatorId(t *testing.T) {
	var validatorId uint64 = 2
	n := structs.Validator{
	}
	var validatorsByValidatorId = make([]structs.Validator, 0)
	validatorsByValidatorId = append(validatorsByValidatorId, n)
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		params     structs.QueryParams
		validators []structs.Validator
		dbResponse error
		code       int
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
			name:   "invalid id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=test",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.QueryParams{
				ValidatorId: validatorId,
			},
			dbResponse: handler.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.QueryParams{
				ValidatorId: validatorId,
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.QueryParams{
				ValidatorId: validatorId,
			},
			validators: validatorsByValidatorId,
			code:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 2 {
				mockDB.EXPECT().GetValidators(tt.req.Context(), tt.params).Return(tt.validators, tt.dbResponse)
			}
			contractor := *handler.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidators)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

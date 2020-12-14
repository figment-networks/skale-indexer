package validator

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client/structs"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGetValidatorByDateRange(t *testing.T) {
	vldByDateRange := structs.Validator{
		Name:        "name_test",
		Description: "description",
	}
	from, _ := time.Parse(handler.Layout, "2006-01-02T15:04:05.000Z")
	to, _ := time.Parse(handler.Layout, "2106-01-02T15:04:05.000Z")
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
			name:   "empty from and to ",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=&to=",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "invalid date from and to ",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2020&to=2100",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 5,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: big.NewInt(0),
			},
			dbResponse: handler.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 6,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: big.NewInt(0),
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 7,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: big.NewInt(0),
			},
			validators: []structs.Validator{vldByDateRange},
			code:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 4 {
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

package api_endpoints

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/client/transport/webapi"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestDelegations(t *testing.T) {
	dlgByDateRange := structs.Delegation{}
	from, _ := time.Parse(structs.Layout, "2006-01-02T15:04:05.000Z")
	to, _ := time.Parse(structs.Layout, "2106-01-02T15:04:05.000Z")
	tests := []struct {
		number      int
		name        string
		req         *http.Request
		params      structs.DelegationParams
		delegations []structs.Delegation
		dbResponse  error
		code        int
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
			name:   "missing parameter all",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "invalid date from and to",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006&to=2106",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "record not found error for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=1903",
				},
			},
			params: structs.DelegationParams{
				DelegationId: "1903",
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=test",
				},
			},
			params: structs.DelegationParams{
				DelegationId: "test",
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=1903",
				},
			},
			params: structs.DelegationParams{
				DelegationId: "1903",
			},
			delegations: []structs.Delegation{dlgByDateRange},
			code:        http.StatusOK,
		},
		{
			number: 7,
			name:   "record not found error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 8,
			name:   "internal server error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 9,
			name:   "success response for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			delegations: []structs.Delegation{dlgByDateRange},
			code:        http.StatusOK,
		},
		{
			number: 10,
			name:   "success response for created time range without validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.DelegationParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			delegations: []structs.Delegation{dlgByDateRange},
			code:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetDelegations(tt.req.Context(), tt.params).Return(tt.delegations, tt.dbResponse)
			}
			contractor := *client.NewClient(mockDB)
			connector := webapi.NewClientConnector(&contractor)
			res := http.HandlerFunc(connector.GetDelegations)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

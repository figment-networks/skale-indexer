package delegations

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

func TestGetValidatorByDateRange(t *testing.T) {
	dlgByDateRange := structs.Delegation{}
	from, _ := time.Parse(structs.Layout, "2006-01-02T15:04:05.000Z")
	to, _ := time.Parse(structs.Layout, "2106-01-02T15:04:05.000Z")
	tests := []struct {
		number      int
		name        string
		req         *http.Request
		params      structs.QueryParams
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
				TimeFrom: from,
				TimeTo:   to,
			},
			dbResponse: structs.ErrNotFound,
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
				TimeFrom: from,
				TimeTo:   to,
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
			if tt.number > 4 {
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
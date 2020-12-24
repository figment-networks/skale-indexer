package webapi

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	storeMocks "github.com/figment-networks/skale-indexer/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	from, _ := time.Parse(structs.Layout, "2006-01-02T15:04:05.000Z")
	to, _ := time.Parse(structs.Layout, "2106-01-02T15:04:05.000Z")
	validatorId := uint64(2)

	tests := []struct {
		name           string
		ttype          string
		req            *http.Request
		expectedParams interface{}
		expectedReturn interface{}
		dbResponse     error
		code           int
	}{

		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "event",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "missing parameters all",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype: "event",
			code:  http.StatusBadRequest,
		},
		{
			name: "bad parameter from and to first check",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006&to=2106",
				},
			},
			ttype: "event",
			code:  http.StatusBadRequest,
		},
		{
			name: "missing parameter id when type is available",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			ttype: "event",
			code:  http.StatusBadRequest,
		},
		{
			name: "missing parameter type when id is available",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			ttype: "event",
			code:  http.StatusBadRequest,
		},
		{
			name: "bad parameter id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator&id=wrong_parameter&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			ttype: "event",
			code:  http.StatusBadRequest,
		},
		{
			name: "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator&id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.EventParams{
				TimeFrom: from,
				TimeTo:   to,
				Id:       validatorId,
				Type:     "validator",
			},
			ttype:      "event",
			dbResponse: structs.ErrNotFound,
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator&id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			dbResponse: errors.New("internal error"),
			expectedParams: structs.EventParams{
				TimeFrom: from,
				TimeTo:   to,
				Id:       validatorId,
				Type:     "validator",
			},
			ttype: "event",
			code:  http.StatusInternalServerError,
		},
		{
			name: "success response for validator",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator&id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.EventParams{
				TimeFrom: from,
				TimeTo:   to,
				Id:       validatorId,
				Type:     "validator",
			},
			expectedReturn: []structs.ContractEvent{{}},
			code:           http.StatusOK,
			ttype:          "event",
		},
		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "delegation",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "missing parameter all",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype: "delegation",
			code:  http.StatusBadRequest,
		},
		{
			name: "invalid date from and to",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006&to=2106",
				},
			},
			ttype: "delegation",
			code:  http.StatusBadRequest,
		},
		{
			name: "record not found error for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationId: "1903",
				TimeFrom:     from,
				TimeTo:       to,
			},
			ttype:      "delegation",
			dbResponse: structs.ErrNotFound,
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=test&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationId: "test",
				TimeFrom:     from,
				TimeTo:       to,
			},
			dbResponse: errors.New("internal error"),
			ttype:      "delegation",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationId: "1903",
				TimeFrom:     from,
				TimeTo:       to,
			},
			ttype:          "delegation",
			expectedReturn: []structs.Delegation{{}},
			code:           http.StatusOK,
		},
		{
			name: "record not found error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:      "delegation",
			dbResponse: structs.ErrNotFound,
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:      "delegation",
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:          "delegation",
			expectedReturn: []structs.Delegation{{}},
			code:           http.StatusOK,
		},
		{
			name: "success response for created time range without validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			ttype:          "delegation",
			expectedReturn: []structs.Delegation{{}},
			code:           http.StatusOK,
		},
		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "delegation_time_line",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "missing parameter all",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype: "delegation_time_line",
			code:  http.StatusBadRequest,
		},
		{
			name: "invalid date from and to",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006&to=2106",
				},
			},
			ttype: "delegation_time_line",
			code:  http.StatusBadRequest,
		},
		{
			name: "record not found error for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationId: "1903",
				TimeFrom:     from,
				TimeTo:       to,
			},
			ttype:      "delegation_time_line",
			dbResponse: structs.ErrNotFound,
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=test&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationId: "test",
				TimeFrom:     from,
				TimeTo:       to,
			},
			dbResponse: errors.New("internal error"),
			ttype:      "delegation_time_line",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for delegation_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationId: "1903",
				TimeFrom:     from,
				TimeTo:       to,
			},
			ttype:          "delegation_time_line",
			expectedReturn: []structs.Delegation{{}},
			code:           http.StatusOK,
		},
		{
			name: "record not found error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:      "delegation_time_line",
			dbResponse: structs.ErrNotFound,
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:      "delegation_time_line",
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorId: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:          "delegation_time_line",
			expectedReturn: []structs.Delegation{{}},
			code:           http.StatusOK,
		},
		{
			name: "success response for created time range without validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			ttype:          "delegation_time_line",
			expectedReturn: []structs.Delegation{{}},
			code:           http.StatusOK,
		},
		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "node",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype:          "node",
			expectedParams: structs.NodeParams{},
			dbResponse:     structs.ErrNotFound,
			code:           http.StatusNotFound,
		},
		{
			name: "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype:          "node",
			expectedParams: structs.NodeParams{},
			dbResponse:     errors.New("internal error"),
			code:           http.StatusInternalServerError,
		},
		{
			name: "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype:          "node",
			expectedParams: structs.NodeParams{},
			expectedReturn: []structs.Node{{}},
			code:           http.StatusOK,
		},
		{
			name: "record not found error with validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				ValidatorId: "2",
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			name: "internal server error with validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				ValidatorId: "2",
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response with validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				ValidatorId: "2",
			},
			expectedReturn: []structs.Node{{}},
			code:           http.StatusOK,
		},
		{
			name: "record not found error with validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&recent=true",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				ValidatorId: "2",
				Recent:      true,
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			name: "internal server error with validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&recent=true",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				ValidatorId: "2",
				Recent:      true,
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response with validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&recent=true",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				ValidatorId: "2",
				Recent:      true,
			},
			expectedReturn: []structs.Node{{}},
			code:           http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.ttype+" - "+tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDB := storeMocks.NewMockDataStore(mockCtrl)
			contractor := *client.NewClient(mockDB)
			connector := NewClientConnector(&contractor)

			if tt.expectedParams != nil {
				switch tt.expectedParams.(type) {
				case structs.DelegationParams:
					if tt.ttype == "delegation" {
						mockDB.EXPECT().GetDelegations(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
					} else {
						mockDB.EXPECT().GetDelegationTimeline(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
					}
				case structs.EventParams:
					mockDB.EXPECT().GetContractEvents(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
				case structs.NodeParams:
					mockDB.EXPECT().GetNodes(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
				}
			}

			var res http.HandlerFunc
			switch tt.ttype {
			case "delegation":
				res = http.HandlerFunc(connector.GetDelegations)
			case "delegation_time_line":
				res = http.HandlerFunc(connector.GetDelegationsTimeline)
			case "event":
				res = http.HandlerFunc(connector.GetContractEvents)
			case "node":
				res = http.HandlerFunc(connector.GetNodes)
			}

			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

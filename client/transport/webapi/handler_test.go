package webapi

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	storeMocks "github.com/figment-networks/skale-indexer/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
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
			name: "success response for created time range without validator_id",
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
			name: "success response for created time range without validator_id",
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
			ttype: "account",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "internal server error for type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator",
				},
			},
			expectedParams: structs.AccountParams{
				Type: "validator",
			},
			dbResponse: errors.New("internal error"),
			ttype:      "account",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator",
				},
			},
			expectedParams: structs.AccountParams{
				Type: "validator",
			},
			ttype:          "account",
			expectedReturn: []structs.Account{},
			code:           http.StatusOK,
		},
		{
			name: "internal server error for address",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "address=0xbeeb437eede0e62a796d9e9c337f62746e925832",
				},
			},
			expectedParams: structs.AccountParams{
				Address: "0xbeeb437eede0e62a796d9e9c337f62746e925832",
			},
			ttype:      "account",
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for address",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "address=0xbeeb437eede0e62a796d9e9c337f62746e925832",
				},
			},
			expectedParams: structs.AccountParams{
				Address: "0xbeeb437eede0e62a796d9e9c337f62746e925832",
			},
			ttype:          "account",
			expectedReturn: []structs.Account{{}},
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
			name: "internal server error with node_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "node_id=2",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				NodeId: "2",
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response with node_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "node_id=2",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				NodeId: "2",
			},
			expectedReturn: []structs.Node{{}},
			code:           http.StatusOK,
		},
		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "validator",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "missing parameter all",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype: "validator",
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
			ttype: "validator",
			code:  http.StatusBadRequest,
		},
		{
			name: "internal server error for validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=test&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				ValidatorId: "test",
				TimeFrom:    from,
				TimeTo:      to,
			},
			dbResponse: errors.New("internal error"),
			ttype:      "validator",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				ValidatorId: "1903",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:          "validator",
			expectedReturn: []structs.Validator{},
			code:           http.StatusOK,
		},
		{
			name: "internal server error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			ttype:      "validator",
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			ttype:          "validator",
			expectedReturn: []structs.Validator{},
			code:           http.StatusOK,
		},
		{
			name: "success response for created time range without validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			ttype:          "validator",
			expectedReturn: []structs.Validator{{}},
			code:           http.StatusOK,
		},
		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "validator_statistics",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "no parameter is valid",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			expectedParams: structs.ValidatorStatisticsParams{},
			dbResponse:     errors.New("internal error"),
			ttype:          "validator_statistics",
			code:           http.StatusInternalServerError,
		},
		{
			name: "internal server error for validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=test",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorId: "test",
			},
			dbResponse: errors.New("internal error"),
			ttype:      "validator_statistics",
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error for validator_id and statistics_type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=1&statistics_type=test",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorId:      "1",
				StatisticsTypeVS: "test",
			},
			dbResponse: errors.New("internal error"),
			ttype:      "validator_statistics",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=1903&statistics_type=1",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorId:      "1903",
				StatisticsTypeVS: "1",
			},
			ttype: "validator_statistics",
			expectedReturn: []structs.ValidatorStatistics{
				{
					ValidatorId: big.NewInt(0),
				},
			},
			code: http.StatusOK,
		},


		{
			name: "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			ttype: "validator_statistics_chart",
			code:  http.StatusMethodNotAllowed,
		},
		{
			name: "missing parameter",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			ttype: "validator_statistics_chart",
			code:  http.StatusBadRequest,
		},
		{
			name: "internal server error for validator_id and statistics_type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=1&statistics_type=test",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorId:      "1",
				StatisticsTypeVS: "test",
			},
			dbResponse: errors.New("internal error"),
			ttype:      "validator_statistics_chart",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=1903&statistics_type=1",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorId:      "1903",
				StatisticsTypeVS: "1",
			},
			ttype: "validator_statistics_chart",
			expectedReturn: []structs.ValidatorStatistics{
				{
					ValidatorId: big.NewInt(0),
				},
			},
			code: http.StatusOK,
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
				case structs.AccountParams:
					mockDB.EXPECT().GetAccounts(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
				case structs.EventParams:
					mockDB.EXPECT().GetContractEvents(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
				case structs.NodeParams:
					mockDB.EXPECT().GetNodes(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
				case structs.ValidatorParams:
					mockDB.EXPECT().GetValidators(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
				case structs.ValidatorStatisticsParams:
					if tt.ttype == "validator_statistics" {
						mockDB.EXPECT().GetValidatorStatistics(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
					} else {
						mockDB.EXPECT().GetValidatorStatisticsChart(tt.req.Context(), tt.expectedParams).Return(tt.expectedReturn, tt.dbResponse)
					}
				}
			}

			var res http.HandlerFunc
			switch tt.ttype {
			case "delegation":
				res = http.HandlerFunc(connector.GetDelegations)
			case "delegation_time_line":
				res = http.HandlerFunc(connector.GetDelegationsTimeline)
			case "account":
				res = http.HandlerFunc(connector.GetAccounts)
			case "event":
				res = http.HandlerFunc(connector.GetContractEvents)
			case "node":
				res = http.HandlerFunc(connector.GetNodes)
			case "validator":
				res = http.HandlerFunc(connector.GetValidators)
			case "validator_statistics":
				res = http.HandlerFunc(connector.GetValidatorStatistics)
			case "validator_statistics_chart":
				res = http.HandlerFunc(connector.GetValidatorStatisticsChart)
			}

			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

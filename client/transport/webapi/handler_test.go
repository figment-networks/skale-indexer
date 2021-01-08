package webapi

import (
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	storeMocks "github.com/figment-networks/skale-indexer/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TODO(lukanus): Add REAL returns

func TestHandler(t *testing.T) {
	from, _ := time.Parse(structs.Layout, "2006-01-02T15:04:05.000Z")
	to, _ := time.Parse(structs.Layout, "2106-01-02T15:04:05.000Z")
	validatorID := uint64(2)

	tests := []struct {
		name             string
		ttype            string
		req              *http.Request
		expectedParams   interface{}
		expectedDBReturn interface{}
		dbResponse       error
		code             int
	}{
		{
			name: "missing parameters all",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			expectedParams: nil,
			ttype:          "event",
			code:           http.StatusBadRequest,
		},
		{
			name: "bad parameter from and to first check",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006&to=2106",
				},
			},
			expectedParams: nil,
			ttype:          "event",
			code:           http.StatusBadRequest,
		},
		{
			name: "missing parameter id when type is available",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator",
				},
			},
			expectedParams: nil,
			ttype:          "event",
			code:           http.StatusBadRequest,
		},
		{
			name: "missing parameter type when id is available",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=2",
				},
			},
			expectedParams: nil,
			ttype:          "event",
			code:           http.StatusBadRequest,
		},
		{
			name: "bad parameter id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "type=validator&id=wrong_parameter",
				},
			},
			expectedParams: nil,
			ttype:          "event",
			code:           http.StatusBadRequest,
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
				Id:       validatorID,
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
				Id:       validatorID,
				Type:     "validator",
			},
			expectedDBReturn: []structs.ContractEvent{{}},
			code:             http.StatusOK,
			ttype:            "event",
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
					RawQuery: "id=test&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationID: "test",
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
					RawQuery: "id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationID: "1903",
				TimeFrom:     from,
				TimeTo:       to,
			},
			ttype:            "delegation",
			expectedDBReturn: []structs.Delegation{{}},
			code:             http.StatusOK,
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
				ValidatorID: "100",
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
				ValidatorID: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:            "delegation",
			expectedDBReturn: []structs.Delegation{{}},
			code:             http.StatusOK,
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
			ttype:            "delegation",
			expectedDBReturn: []structs.Delegation{{}},
			code:             http.StatusOK,
		},
		{
			name: "invalid date from and to",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006&to=2106&timeline=1",
				},
			},
			expectedParams: nil,
			ttype:          "delegation",
			code:           http.StatusBadRequest,
		},
		{
			name: "internal server error for delegation id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=test&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z&timeline=1",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationID: "test",
				TimeFrom:     from,
				TimeTo:       to,
			},
			dbResponse: errors.New("internal error"),
			ttype:      "delegation",
			code:       http.StatusInternalServerError,
		},
		{
			name: "success response for delegation id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z&timeline=1",
				},
			},
			expectedParams: structs.DelegationParams{
				DelegationID: "1903",
				TimeFrom:     from,
				TimeTo:       to,
			},
			ttype:            "delegation",
			expectedDBReturn: []structs.Delegation{{}},
			code:             http.StatusOK,
		},
		{
			name: "internal server error for created time range",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z&timeline=1",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorID: "100",
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
					RawQuery: "validator_id=100&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z&timeline=1",
				},
			},
			expectedParams: structs.DelegationParams{
				ValidatorID: "100",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:            "delegation",
			expectedDBReturn: []structs.Delegation{{}},
			code:             http.StatusOK,
		},
		{
			name: "success response for created time range without validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z&timeline=1",
				},
			},
			expectedParams: structs.DelegationParams{
				TimeFrom: from,
				TimeTo:   to,
			},
			ttype:            "delegation",
			expectedDBReturn: []structs.Delegation{{}},
			code:             http.StatusOK,
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
			ttype:            "account",
			expectedDBReturn: []structs.Account{{}},
			code:             http.StatusOK,
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
			ttype:            "account",
			expectedDBReturn: []structs.Account{{}},
			code:             http.StatusOK,
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
			ttype:            "node",
			expectedParams:   structs.NodeParams{},
			expectedDBReturn: []structs.Node{{}},
			code:             http.StatusOK,
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
			expectedDBReturn: []structs.Node{{}},
			code:             http.StatusOK,
		},
		{
			name: "internal server error with node id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=2",
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
			name: "success response with node id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=2",
				},
			},
			ttype: "node",
			expectedParams: structs.NodeParams{
				NodeId: "2",
			},
			expectedDBReturn: []structs.Node{{}},
			code:             http.StatusOK,
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
					RawQuery: "id=test&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				ValidatorID: "test",
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
					RawQuery: "id=1903&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			expectedParams: structs.ValidatorParams{
				ValidatorID: "1903",
				TimeFrom:    from,
				TimeTo:      to,
			},
			ttype:            "validator",
			expectedDBReturn: []structs.Validator{{}},
			code:             http.StatusOK,
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
			ttype:            "validator",
			expectedDBReturn: []structs.Validator{{}},
			code:             http.StatusOK,
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
			ttype:            "validator",
			expectedDBReturn: []structs.Validator{{}},
			code:             http.StatusOK,
		},
		{
			name: "parameters are invalid",
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
					RawQuery: "id=123",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorID: "123",
			},
			dbResponse: errors.New("internal error"),
			ttype:      "validator_statistics",
			code:       http.StatusInternalServerError,
		},
		{
			name: "internal server error for validator_id and type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=1903&type=FEE",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorID: "1903",
				Type:        structs.ValidatorStatisticsTypeFee,
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
					RawQuery: "id=1903&type=FEE",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorID: "1903",
				Type:        structs.ValidatorStatisticsTypeFee,
			},
			ttype:            "validator_statistics",
			expectedDBReturn: []structs.ValidatorStatistics{{ValidatorID: big.NewInt(1903)}},
			code:             http.StatusOK,
		},
		{
			name: "internal server error for validator_id and type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=1&type=FEE&timeline=true",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorID: "1",
				Type:        structs.ValidatorStatisticsTypeFee,
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
					RawQuery: "id=1903&type=FEE&timeline=1",
				},
			},
			expectedParams: structs.ValidatorStatisticsParams{
				ValidatorID: "1903",
				Type:        structs.ValidatorStatisticsTypeFee,
			},
			ttype:            "validator_statistics",
			expectedDBReturn: []structs.ValidatorStatistics{{Type: structs.ValidatorStatisticsTypeFee, ValidatorID: big.NewInt(1903)}},
			code:             http.StatusOK,
		},
	}

	for _, tt := range tests {
		//	tt := tt
		t.Run(tt.ttype+" - "+tt.name, func(t *testing.T) {
			//		t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDB := storeMocks.NewMockDataStore(mockCtrl)
			zl := zaptest.NewLogger(t)
			contractor := *client.NewClient(zl, mockDB)
			connector := NewClientConnector(&contractor)

			if tt.expectedParams != nil {
				switch tt.expectedParams.(type) {
				case structs.DelegationParams:
					if strings.Contains(tt.req.URL.RawQuery, "timeline") {
						mockDB.EXPECT().GetDelegationTimeline(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
					} else {
						mockDB.EXPECT().GetDelegations(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
					}
				case structs.AccountParams:
					mockDB.EXPECT().GetAccounts(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
				case structs.EventParams:
					mockDB.EXPECT().GetContractEvents(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
				case structs.NodeParams:
					mockDB.EXPECT().GetNodes(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
				case structs.ValidatorParams:
					mockDB.EXPECT().GetValidators(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
				case structs.ValidatorStatisticsParams:
					if strings.Contains(tt.req.URL.RawQuery, "timeline") {
						mockDB.EXPECT().GetValidatorStatisticsTimeline(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
					} else {
						mockDB.EXPECT().GetValidatorStatistics(tt.req.Context(), tt.expectedParams).Return(tt.expectedDBReturn, tt.dbResponse)
					}
				}
			}

			var res http.HandlerFunc
			switch tt.ttype {
			case "delegation":
				res = http.HandlerFunc(connector.GetDelegation)
			case "account":
				res = http.HandlerFunc(connector.GetAccount)
			case "event":
				res = http.HandlerFunc(connector.GetContractEvents)
			case "node":
				res = http.HandlerFunc(connector.GetNode)
			case "validator":
				res = http.HandlerFunc(connector.GetValidator)
			case "validator_statistics":
				res = http.HandlerFunc(connector.GetValidatorStatistics)
			}

			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			require.Equal(t, rr.Code, tt.code)

		})
	}
}

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
// TODO: change failed tests after handler changes
func TestGetAllEvents(t *testing.T) {
	from, _ := time.Parse(structs.Layout, "2006-01-02T15:04:05.000Z")
	to, _ := time.Parse(structs.Layout, "2106-01-02T15:04:05.000Z")
	var validatorId uint64 = 2
	var events = make([]structs.ContractEvent, 0)
	events = append(events, structs.ContractEvent{})
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		params     structs.QueryParams
		events     []structs.ContractEvent
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
			name:   "missing parameters all",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "bad parameter from and to first check",
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
			name:   "missing parameter bound_type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 5,
			name:   "bad parameter bound_type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=unsupported&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 6,
			name:   "missing parameter validator_id for validator",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=validator&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 7,
			name:   "bad parameter validator_id for validator",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=validator&validator_id=wrong_parameter&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 8,
			name:   "missing parameter validator_id for delegation",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=delegation&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 9,
			name:   "bad parameter validator_id for delegation",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=delegation&validator_id=wrong_parameter&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 10,
			name:   "bad parameter delegation_id for delegation",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=delegation&validator_id=2&delegation_id=wrong_type&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 11,
			name:   "record not found error for validator",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=validator&validator_id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: validatorId,
				BoundType:   "validator",
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 12,
			name:   "internal server error for validator",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=validator&validator_id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			dbResponse: errors.New("internal error"),
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: validatorId,
				BoundType:   "validator",
			},
			code: http.StatusInternalServerError,
		},
		{
			number: 13,
			name:   "success response for validator",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=validator&validator_id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: validatorId,
				BoundType:   "validator",
			},
			events: events,
			code:   http.StatusOK,
		},
		{
			number: 14,
			name:   "record not found error for delegation",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=delegation&validator_id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: validatorId,
				BoundType:   "delegation",
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 15,
			name:   "internal server error for delegation",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=delegation&validator_id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			dbResponse: errors.New("internal error"),
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: validatorId,
				BoundType:   "delegation",
			},
			code: http.StatusInternalServerError,
		},
		{
			number: 16,
			name:   "success response for delegation",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "bound_type=delegation&validator_id=2&from=2006-01-02T15:04:05.000Z&to=2106-01-02T15:04:05.000Z",
				},
			},
			params: structs.QueryParams{
				TimeFrom:    from,
				TimeTo:      to,
				ValidatorId: validatorId,
				BoundType:   "delegation",
			},
			events: events,
			code:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 10 {
				mockDB.EXPECT().GetContractEvents(tt.req.Context(), tt.params).Return(tt.events, tt.dbResponse)
			}
			contractor := *client.NewClient(mockDB)
			connector := webapi.NewClientConnector(&contractor)
			res := http.HandlerFunc(connector.GetContractEvents)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)

			// TODO: check response body
		})
	}
}

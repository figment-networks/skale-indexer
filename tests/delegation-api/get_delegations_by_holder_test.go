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
	"time"
)

var dlgsByHolder = make([]structs.Delegation, 1)

func TestGetDelegationsByHolder(t *testing.T) {
	var holder uint64 = 1
	dlg := structs.Delegation{
		Holder:               holder,
		ValidatorId:          uint64(2),
		Amount:               uint64(0),
		DelegationPeriod:     uint64(0),
		Created:              time.Now(),
		Started:              time.Now(),
		Finished:             time.Now(),
		Info:                 "info1",
		Status:               1,
		SmartContractIndex:   1903,
		SmartContractAddress: 1001,
	}
	dlgsByHolder = append(dlgsByHolder, dlg)
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
			name:   "empty parameter",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "holder=",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "invalid id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "holder=test",
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
					RawQuery: "holder=1",
				},
			},
			params: structs.QueryParams{
				Holder: holder,
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
					RawQuery: "holder=1",
				},
			},
			params: structs.QueryParams{
				Holder: holder,
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
					RawQuery: "holder=1",
				},
			},
			params: structs.QueryParams{
				Holder: holder,
			},
			delegations: dlgsByHolder,
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
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetDelegations)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

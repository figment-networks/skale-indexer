package tests

import (
	"../../client"
	"../../handler"
	"../../store"
	"../../structs"
	"errors"
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
	holder := "holder1"
	var validatorId uint64 = 2
	var amount uint64 = 0
	var delegationPeriod uint64 = 0
	var created time.Time = time.Now()
	var started time.Time = time.Now()
	var finished time.Time = time.Now()
	info := "info1"
	dlg := structs.Delegation{
		Holder:           &holder,
		ValidatorId:      &validatorId,
		Amount:           &amount,
		DelegationPeriod: &delegationPeriod,
		Created:          &created,
		Started:          &started,
		Finished:         &finished,
		Info:             &info,
	}
	dlgsByHolder = append(dlgsByHolder, dlg)
	tests := []struct {
		number      int
		name        string
		req         *http.Request
		holder      *string
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
			name:   "bad request",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "bad request",
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
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "holder=holder1",
				},
			},
			holder:     &holder,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "holder=holder1",
				},
			},
			holder:      &holder,
			delegations: dlgsByHolder,
			code:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number == 4 || tt.number == 5 {
				mockDB.EXPECT().GetDelegationsByHolder(tt.req.Context(), tt.holder).Return(tt.delegations, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetDelegationsByHolder)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

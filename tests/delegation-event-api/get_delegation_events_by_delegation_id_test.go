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

var dlgEventsByDelegationId = make([]structs.DelegationEvent, 1)

func TestGetDelegationEventsByHolder(t *testing.T) {
	invalidId := "delegationId1"
	id := "11053aa6-4bbb-4094-b588-8368cd621f2c"
	eventName := "eventName1"
	eventTime := time.Now()
	dlg := structs.DelegationEvent{
		DelegationId: id,
		EventName:    eventName,
		EventTime:    eventTime,
	}
	dlgEventsByDelegationId = append(dlgEventsByDelegationId, dlg)
	tests := []struct {
		number           int
		name             string
		req              *http.Request
		delegationId     string
		delegationEvents []structs.DelegationEvent
		dbResponse       error
		code             int
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
			name:   "empty id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=11053aa6-4bbb-4094-b588-8368cd621f2c",
				},
			},
			delegationId: id,
			dbResponse:   errors.New("record not found"),
			code:         http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=delegationId1",
				},
			},
			delegationId: invalidId,
			dbResponse:   errors.New("internal error"),
			code:         http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=11053aa6-4bbb-4094-b588-8368cd621f2c",
				},
			},
			delegationId:     id,
			delegationEvents: dlgEventsByDelegationId,
			code:             http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetDelegationEventsByDelegationId(tt.req.Context(), tt.delegationId).Return(tt.delegationEvents, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetDelegationEventsByDelegationId)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

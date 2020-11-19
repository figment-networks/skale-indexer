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
	delegationId := "delegationId1"
	eventName := "eventName1"
	eventTime := time.Now()
	dlg := structs.DelegationEvent{
		DelegationId: &delegationId,
		EventName:    &eventName,
		EventTime:    &eventTime,
	}
	dlgEventsByDelegationId = append(dlgEventsByDelegationId, dlg)
	tests := []struct {
		number           int
		name             string
		req              *http.Request
		delegationId     *string
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
					RawQuery: "delegation_id=delegationId1",
				},
			},
			delegationId: &delegationId,
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
			delegationId: &delegationId,
			dbResponse:   errors.New("internal error"),
			code:         http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "delegation_id=delegationId1",
				},
			},
			delegationId:     &delegationId,
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

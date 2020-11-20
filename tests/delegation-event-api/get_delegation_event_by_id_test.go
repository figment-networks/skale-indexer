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

var dlgById structs.DelegationEvent

func TestGetDelegationEventById(t *testing.T) {
	delegationId := "delegationId1"
	eventName := "eventName1"
	eventTime := time.Now()

	dlgById = structs.DelegationEvent{
		DelegationId: delegationId,
		EventName:    eventName,
		EventTime:    eventTime,
	}
	var id = "id_test"
	tests := []struct {
		number          int
		name            string
		req             *http.Request
		id              string
		delegationEvent structs.DelegationEvent
		dbResponse      error
		code            int
	}{
		{
			number: 1,
			name:   "not allowed method",
			req: &http.Request{
				Method: http.MethodPost,
			},
			id:   id,
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
			name:   "invalid id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=",
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
					RawQuery: "id=id_test",
				},
			},
			id:         id,
			dbResponse: errors.New("record not found"),
			code:       http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=id_test",
				},
			},
			id:         id,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=id_test",
				},
			},
			id:              id,
			delegationEvent: dlgById,
			code:            http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetDelegationEventById(tt.req.Context(), tt.id).Return(tt.delegationEvent, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetDelegationEventById)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

package tests

import (
	"bytes"
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	invalidSyntaxForDelegationEvents = `[{
        "delegation_id": "delegation_id_test",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z",
	`
	invalidPropertyNameForDelegationEvents = `[{
    	"delegation_id_invalid": "delegation_id_test",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
	}]`
	validJsonForDelegationEvents = `[{
		"delegation_id": "122bcb7c-b283-4f59-a945-75d6cf37e978",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
    },	
	{
    	"delegation_id": "11053aa6-4bbb-4094-b588-8368cd621f2c",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
    }
	]`
	// same value should be used in json examples above for valid cases
	dummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleDelegations []structs.DelegationEvent

func TestSaveOrUpdateDelegations(t *testing.T) {
	delegationId := "122bcb7c-b283-4f59-a945-75d6cf37e978"
	eventName := "event_name_test"
	layout := "2006-01-02T15:04:05.000Z"
	exampleTime, _ := time.Parse(layout, dummyTime)
	var eventTime = exampleTime
	example1 := structs.DelegationEvent{
		DelegationId: delegationId,
		EventName:    eventName,
		EventTime:    eventTime,
	}
	delegationId2 := "11053aa6-4bbb-4094-b588-8368cd621f2c"
	eventName2 := "event_name_test"
	example2 := structs.DelegationEvent{
		DelegationId: delegationId2,
		EventName:    eventName2,
		EventTime:    eventTime,
	}
	exampleDelegations = append(exampleDelegations, example1)
	exampleDelegations = append(exampleDelegations, example2)

	tests := []struct {
		number           int
		name             string
		req              *http.Request
		delegationEvents []structs.DelegationEvent
		dbResponse       error
		code             int
	}{
		{
			number: 1,
			name:   "not allowed method",
			req: &http.Request{
				Method: http.MethodGet,
			},
			code: http.StatusMethodNotAllowed,
		},
		{
			number: 2,
			name:   "invalid json syntax request body",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForDelegationEvents))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "missing parameter",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForDelegationEvents))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForDelegationEvents))),
			},
			delegationEvents: exampleDelegations,
			dbResponse:       errors.New("internal error"),
			code:             http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForDelegationEvents))),
			},
			delegationEvents: exampleDelegations,
			code:             http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().SaveOrUpdateDelegationEvents(tt.req.Context(), tt.delegationEvents).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateDelegationEvents)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

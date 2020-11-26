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
	invalidSyntaxForEvents = `[{
        "delegation_id": "delegation_id_test",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z",
	`
	invalidPropertyNameForEvents = `[{
    	"delegation_id_invalid": "delegation_id_test",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
	}]`
	validJsonForEvents = `[{
		"block_height": 100,
		"smart_contract_address": "smart_contract_address1",
		"transaction_index": 15,
		"event_type": "eventType1",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z",
	    "event_info": {
						"wallet": "wallet_test1",
						"holder": "holder3",	
						"destination": [10,21],
						"validator_id": 21,
						"amount": 22
				}
    },	
	{
		"block_height": 101,
		"smart_contract_address": "smart_contract_address2",
		"transaction_index": 25,
		"event_type": "eventType2",        
		"event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z",
	    "event_info": {
						"wallet": "wallet_test2",
						"holder": "holder1",	
						"destination":  [1,2],
						"validator_id": 1,
						"amount": 2
				}
    }
	]`
	// same value should be used in json examples above for valid cases
	dummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleEvents []structs.Event

func TestSaveOrUpdateEvents(t *testing.T) {
	blockHeight := int64(100)
	smartContractAddress := "smart_contract_address1"
	transactionIndex := int64(15)
	eventType := "eventType1"
	eventName := "event_name_test"
	exampleTime, _ := time.Parse(handler.Layout, dummyTime)
	var eventTime = exampleTime
	example1 := structs.Event{
		BlockHeight:          blockHeight,
		SmartContractAddress: smartContractAddress,
		TransactionIndex:     transactionIndex,
		EventType:            eventType,
		EventName:            eventName,
		EventTime:            eventTime,
		EventInfo: structs.EventInfo{
			Wallet:      "wallet_test1",
			Holder:      "holder3",
			Destination: []structs.Address{10, 21},
			ValidatorId: 21,
			Amount:      22,
		},
	}
	blockHeight2 := int64(101)
	smartContractAddress2 := "smart_contract_address2"
	transactionIndex2 := int64(25)
	eventType2 := "eventType2"
	eventName2 := "event_name_test"
	example2 := structs.Event{
		BlockHeight:          blockHeight2,
		SmartContractAddress: smartContractAddress2,
		TransactionIndex:     transactionIndex2,
		EventType:            eventType2,
		EventName:            eventName2,
		EventTime:            eventTime,
		EventInfo: structs.EventInfo{
			Wallet:      "wallet_test2",
			Holder:      "holder1",
			Destination: []structs.Address{1, 2},
			ValidatorId: 1,
			Amount:      2,
		},
	}
	exampleEvents = append(exampleEvents, example1)
	exampleEvents = append(exampleEvents, example2)

	tests := []struct {
		number     int
		name       string
		req        *http.Request
		events     []structs.Event
		dbResponse error
		code       int
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
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForEvents))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "missing parameter",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForEvents))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForEvents))),
			},
			events:     exampleEvents,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForEvents))),
			},
			events: exampleEvents,
			code:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().SaveOrUpdateEvents(tt.req.Context(), tt.events).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateEvents)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

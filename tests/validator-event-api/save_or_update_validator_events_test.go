package tests

import (
	"../../client"
	"../../handler"
	"../../store"
	"../../structs"
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	invalidSyntaxForValidatorEvents = `[{
        "validator_id": "validator_id_test",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z",
	`
	invalidPropertyNameForValidatorEvents = `[{
    	"validator_id_invalid": "validator_id_test",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
	}]`
	validJsonForValidatorEvents = `[{
		"validator_id": "validator_id_test1",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
    },	
	{
    	"validator_id": "validator_id_test2",
        "event_name": "event_name_test",
        "event_time": "2014-11-12T11:45:26.371Z"
    }
	]`
	// same value should be used in json examples above for valid cases
	dummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleValidators []structs.ValidatorEvent

func TestSaveOrUpdateValidators(t *testing.T) {
	validatorId := "validator_id_test1"
	eventName := "event_name_test"
	layout := "2006-01-02T15:04:05.000Z"
	exampleTime, _ := time.Parse(layout, dummyTime)
	var eventTime = exampleTime
	example1 := structs.ValidatorEvent{
		ValidatorId: validatorId,
		EventName:   eventName,
		EventTime:   eventTime,
	}
	validatorId2 := "validator_id_test2"
	eventName2 := "event_name_test"
	example2 := structs.ValidatorEvent{
		ValidatorId: validatorId2,
		EventName:   eventName2,
		EventTime:   eventTime,
	}
	exampleValidators = append(exampleValidators, example1)
	exampleValidators = append(exampleValidators, example2)

	tests := []struct {
		number          int
		name            string
		req             *http.Request
		validatorEvents []structs.ValidatorEvent
		dbResponse      error
		code            int
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
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForValidatorEvents))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "bad request",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForValidatorEvents))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForValidatorEvents))),
			},
			validatorEvents: exampleValidators,
			dbResponse:      errors.New("internal error"),
			code:            http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForValidatorEvents))),
			},
			validatorEvents: exampleValidators,
			code:            http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().SaveOrUpdateValidatorEvents(tt.req.Context(), tt.validatorEvents).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateValidatorEvents)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

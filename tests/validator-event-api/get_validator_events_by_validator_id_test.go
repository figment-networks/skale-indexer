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

var vldEventsByValidatorId = make([]structs.ValidatorEvent, 1)

func TestGetValidatorEventsByValidatorId(t *testing.T) {
	validatorId := "validatorId1"
	eventName := "eventName1"
	eventTime := time.Now()
	ve := structs.ValidatorEvent{
		ValidatorId: &validatorId,
		EventName:   &eventName,
		EventTime:   &eventTime,
	}
	vldEventsByValidatorId = append(vldEventsByValidatorId, ve)
	tests := []struct {
		number          int
		name            string
		req             *http.Request
		validatorId     *string
		validatorEvents []structs.ValidatorEvent
		dbResponse      error
		code            int
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
					RawQuery: "validator_id=",
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
					RawQuery: "validator_id=validatorId1",
				},
			},
			validatorId: &validatorId,
			dbResponse:  errors.New("record not found"),
			code:        http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=validatorId1",
				},
			},
			validatorId: &validatorId,
			dbResponse:  errors.New("internal error"),
			code:        http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=validatorId1",
				},
			},
			validatorId:     &validatorId,
			validatorEvents: vldEventsByValidatorId,
			code:            http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetValidatorEventsByValidatorId(tt.req.Context(), tt.validatorId).Return(tt.validatorEvents, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorEventsByValidatorId)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}
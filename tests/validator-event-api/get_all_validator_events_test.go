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

var vldEvents = make([]structs.ValidatorEvent, 1)

func TestGetAllValidatorEvents(t *testing.T) {
	validatorId := "validatorId1"
	eventName := "eventName1"
	eventTime := time.Now()
	vld := structs.ValidatorEvent{
		ValidatorId: validatorId,
		EventName:   eventName,
		EventTime:   eventTime,
	}
	vldEvents = append(vldEvents, vld)
	tests := []struct {
		number          int
		name            string
		req             *http.Request
		validatorId     string
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
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			validatorId: validatorId,
			dbResponse:  errors.New("record not found"),
			code:        http.StatusNotFound,
		},
		{
			number: 3,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
			},
			validatorId: validatorId,
			dbResponse:  errors.New("internal error"),
			code:        http.StatusInternalServerError,
		},
		{
			number: 4,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			validatorId:     validatorId,
			validatorEvents: vldEvents,
			code:            http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 1 {
				mockDB.EXPECT().GetAllValidatorEvents(tt.req.Context()).Return(tt.validatorEvents, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorEvents)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

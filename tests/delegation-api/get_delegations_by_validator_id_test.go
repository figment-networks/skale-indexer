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

var dlgsByValidatorId = make([]structs.Delegation, 1)

func TestGetDelegationsByValidatorId(t *testing.T) {
	holder := "holder1"
	var validatorId uint64 = 2
	var amount uint64 = 0
	var delegationPeriod uint64 = 0
	var created = time.Now()
	var started = time.Now()
	var finished = time.Now()
	info := "info1"
	dlg := structs.Delegation{
		Holder:           holder,
		ValidatorId:      validatorId,
		Amount:           amount,
		DelegationPeriod: delegationPeriod,
		Created:          created,
		Started:          started,
		Finished:         finished,
		Info:             info,
	}
	dlgsByValidatorId = append(dlgsByValidatorId, dlg)
	tests := []struct {
		number      int
		name        string
		req         *http.Request
		validatorId uint64
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
			name:   "invalid id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=test",
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
					RawQuery: "validator_id=2",
				},
			},
			validatorId: validatorId,
			dbResponse:  errors.New("record not found"),
			code:        http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			validatorId: validatorId,
			dbResponse:  errors.New("internal error"),
			code:        http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			validatorId: validatorId,
			delegations: dlgsByValidatorId,
			code:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetDelegationsByValidatorId(tt.req.Context(), tt.validatorId).Return(tt.delegations, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetDelegationsByValidatorId)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

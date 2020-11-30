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

func TestGetDelegationNextEpochStatisticsByValidatorId(t *testing.T) {
	var validatorId uint64 = 2
	s := structs.DelegationStatistics{
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		Status:        1,
		ValidatorId:   2,
		Amount:        3,
		StatisticType: structs.NextEpochStatisticsType,
	}
	var statsByValidatorId = make([]structs.DelegationStatistics, 0)
	statsByValidatorId = append(statsByValidatorId, s)
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		params     structs.QueryParams
		stats      []structs.DelegationStatistics
		dbResponse error
		code       int
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
			name:   "invalid id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=test&statistic_type=next-epoch",
				},
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&statistic_type=next-epoch",
				},
			},
			params: structs.QueryParams{
				ValidatorId:   validatorId,
				StatisticType: structs.NextEpochStatisticsType,
			},
			dbResponse: handler.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&statistic_type=next-epoch",
				},
			},
			params: structs.QueryParams{
				ValidatorId:   validatorId,
				StatisticType: structs.NextEpochStatisticsType,
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&statistic_type=next-epoch",
				},
			},
			params: structs.QueryParams{
				ValidatorId:   validatorId,
				StatisticType: structs.NextEpochStatisticsType,
			},
			stats: statsByValidatorId,
			code:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 2 {
				mockDB.EXPECT().GetDelegationStatistics(tt.req.Context(), tt.params).Return(tt.stats, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetDelegationStatistics)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
			for _, s := range tt.stats {
				assert.True(t, s.StatisticType == structs.NextEpochStatisticsType)
			}
		})
	}
}

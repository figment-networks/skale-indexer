package validator_statistics

import (
	"errors"
	"github.com/figment-networks/skale-indexer/api/structs"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGetAllValidatorActiveNodesStatistics(t *testing.T) {
	d := structs.ValidatorStatistics{
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		ValidatorId:    2,
		Amount:         3,
		ETHBlockHeight: 1000,
		StatisticType:  structs.ValidatorStatisticsTypeActiveNodes,
	}
	var stats = make([]structs.ValidatorStatistics, 0)
	stats = append(stats, d)
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		params     structs.QueryParams
		stats      []structs.ValidatorStatistics
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
			name:   "unknown type",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "statistic_type=unkown",
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
					RawQuery: "statistic_type=active_nodes",
				},
			},
			params: structs.QueryParams{
				StatisticTypeVS: structs.ValidatorStatisticsTypeActiveNodes,
			},
			dbResponse: handler.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 5,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "statistic_type=active_nodes",
				},
			},
			dbResponse: errors.New("internal error"),
			params: structs.QueryParams{
				StatisticTypeVS: structs.ValidatorStatisticsTypeActiveNodes,
			},
			code: http.StatusInternalServerError,
		},
		{
			number: 6,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "statistic_type=active_nodes",
				},
			},
			params: structs.QueryParams{
				StatisticTypeVS: structs.ValidatorStatisticsTypeActiveNodes,
			},
			stats: stats,
			code:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().GetValidatorStatistics(tt.req.Context(), tt.params).Return(tt.stats, tt.dbResponse)
			}
			contractor := *handler.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorStatistics)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
			for _, s := range tt.stats {
				assert.True(t, s.StatisticType == structs.ValidatorStatisticsTypeActiveNodes)
			}
		})
	}
}

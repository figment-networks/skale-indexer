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
)

func TestGetValidatorTotalStakeStatisticsById(t *testing.T) {
	statById := structs.ValidatorStatistics{
		StatisticType: structs.TotalStakeStatisticsTypeVS,
	}
	var id = "11053aa6-4bbb-4094-b588-8368cd621f2c"
	var invalidId = "id_test"
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
			params: structs.QueryParams{
				Id: id,
			},
			code: http.StatusMethodNotAllowed,
		},
		{
			number: 2,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=11053aa6-4bbb-4094-b588-8368cd621f2c&statistic_type=total_stake",
				},
			},
			params: structs.QueryParams{
				Id:              id,
				StatisticTypeVS: structs.TotalStakeStatisticsTypeVS,
			},
			dbResponse: handler.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 3,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=id_test&statistic_type=total_stake",
				},
			},
			params: structs.QueryParams{
				Id:              invalidId,
				StatisticTypeVS: structs.TotalStakeStatisticsTypeVS,
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 4,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "id=11053aa6-4bbb-4094-b588-8368cd621f2c&statistic_type=total_stake",
				},
			},
			params: structs.QueryParams{
				Id:              id,
				StatisticTypeVS: structs.TotalStakeStatisticsTypeVS,
			},
			stats: []structs.ValidatorStatistics{statById},
			code:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 1 {
				mockDB.EXPECT().GetValidatorStatistics(tt.req.Context(), tt.params).Return(tt.stats, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetValidatorStatistics)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
			for _, s := range tt.stats {
				assert.True(t, s.StatisticType == structs.TotalStakeStatisticsTypeVS)
			}
		})
	}
}
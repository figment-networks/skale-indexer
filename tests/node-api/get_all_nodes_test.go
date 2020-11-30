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

func TestGetAllNodes(t *testing.T) {
	n := structs.Node{
		CreatedAt:                time.Time{},
		UpdatedAt:                time.Time{},
		Address:                  uint64(1),
		Name:                     "name1",
		Ip:                       "127.0.0.1",
		PublicIp:                 "127.0.0.1",
		Port:                     8080,
		PublicKey:                "public key",
		StartBlock:               1000,
		LastRewardDate:           time.Now(),
		FinishTime:               time.Now(),
		Status:                   "",
		ValidatorId:              2,
		RegistrationDate:         time.Now(),
		LastBountyCall:           time.Now(),
		CalledGetBountyThisEpoch: true,
		Balance:                  0.1234,
	}
	var nodes = make([]structs.Node, 0)
	nodes = append(nodes, n)
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		params     structs.QueryParams
		nodes      []structs.Node
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
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			params:     structs.QueryParams{},
			dbResponse: handler.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 3,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			dbResponse: errors.New("internal error"),
			params:     structs.QueryParams{},
			code:       http.StatusInternalServerError,
		},
		{
			number: 4,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			params: structs.QueryParams{},
			nodes:  nodes,
			code:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 1 {
				mockDB.EXPECT().GetNodes(tt.req.Context(), tt.params).Return(tt.nodes, tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetNodes)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

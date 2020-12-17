package api_endpoints

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/client/transport/webapi"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNodes(t *testing.T) {
	var nodesByValidatorId = make([]structs.Node, 0)
	nodesByValidatorId = append(nodesByValidatorId, structs.Node{})
	tests := []struct {
		number     int
		name       string
		req        *http.Request
		params     structs.NodeParams
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
			params:     structs.NodeParams{},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 3,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			params:     structs.NodeParams{},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 4,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{},
			},
			params: structs.NodeParams{},
			nodes:  nodesByValidatorId,
			code:   http.StatusOK,
		},
		{
			number: 5,
			name:   "record not found error with validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.NodeParams{
				ValidatorId: "2",
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 6,
			name:   "internal server error with validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.NodeParams{
				ValidatorId: "2",
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 7,
			name:   "success response with validator_id",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.NodeParams{
				ValidatorId: "2",
			},
			nodes: nodesByValidatorId,
			code:  http.StatusOK,
		},
		{
			number: 8,
			name:   "record not found error with validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&recent=true",
				},
			},
			params: structs.NodeParams{
				ValidatorId: "2",
				Recent:      true,
			},
			dbResponse: structs.ErrNotFound,
			code:       http.StatusNotFound,
		},
		{
			number: 9,
			name:   "internal server error with validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&recent=true",
				},
			},
			params: structs.NodeParams{
				ValidatorId: "2",
				Recent:      true,
			},
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 10,
			name:   "success response with validator_id and recent",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2&recent=true",
				},
			},
			params: structs.NodeParams{
				ValidatorId: "2",
				Recent:      true,
			},
			nodes: nodesByValidatorId,
			code:  http.StatusOK,
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
			contractor := *client.NewClient(mockDB)
			connector := webapi.NewClientConnector(&contractor)
			res := http.HandlerFunc(connector.GetNodes)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

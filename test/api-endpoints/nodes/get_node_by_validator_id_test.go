package nodes

import (
	"errors"
	"github.com/figment-networks/skale-indexer/client/structs"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetNodesByValidatorId(t *testing.T) {
	var validatorId int64 = 2
	n := structs.Node{
	}
	var nodesByValidatorId = make([]structs.Node, 0)
	nodesByValidatorId = append(nodesByValidatorId, n)
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
			number: 3,
			name:   "record not found error",
			req: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					RawQuery: "validator_id=2",
				},
			},
			params: structs.QueryParams{
				ValidatorId: big.NewInt(validatorId),
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
					RawQuery: "validator_id=2",
				},
			},
			params: structs.QueryParams{
				ValidatorId: big.NewInt(validatorId),
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
					RawQuery: "validator_id=2",
				},
			},
			params: structs.QueryParams{
				ValidatorId: big.NewInt(validatorId),
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
			if tt.number > 2 {
				mockDB.EXPECT().GetNodes(tt.req.Context(), tt.params).Return(tt.nodes, tt.dbResponse)
			}
			contractor := *handler.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.GetNodes)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

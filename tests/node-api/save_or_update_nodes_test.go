package tests

import (
	"bytes"
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/handler"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	invalidSyntaxForNodes = `[{
        "name": "name1",
        "ip": "127.0.0.1",
	`
	invalidPropertyNameForNodes = `[{
 		"name": "name1",
        "ip": "127.0.0.1",
        "public_ip": "127.0.0.1",
        "port": "1903",
        "public_key": "public key1",
        "start_block": "start block1",
        "last_reward_date": "2014-11-12T11:45:26.371Z",
        "finish_time": "2014-11-12T11:45:26.371Z",
        "status": "Active",
        "validator_id_invalid": 2
	}]`
	validJsonForNodes = `[{
 		"name": "name1",
        "ip": "127.0.0.1",
        "public_ip": "127.0.0.1",
        "port": 1903,
        "public_key": "public key1",
        "start_block": 1000,
        "last_reward_date": "2014-11-12T11:45:26.371Z",
        "finish_time": "2014-11-12T11:45:26.371Z",
        "status": "Active",
        "validator_id": 2
	},	
	{
 		"name": "name2",
        "ip": "127.0.0.2",
        "public_ip": "127.0.0.2",
        "port": 1904,
        "public_key": "public key2",
        "start_block": 1001,
        "last_reward_date": "2014-11-12T11:45:26.371Z",
        "finish_time": "2014-11-12T11:45:26.371Z",
        "status": "Pending",
        "validator_id": 3
	}
	]`
	// same value should be used in json examples above for valid cases
	dummyTime = "2014-11-12T11:45:26.371Z"
)

var exampleNodes []structs.Node

func TestSaveOrUpdateNodes(t *testing.T) {
	exampleTime, _ := time.Parse(handler.Layout, dummyTime)
	example1 := structs.Node{
		Name:           "name1",
		Ip:             "127.0.0.1",
		PublicIp:       "127.0.0.1",
		Port:           1903,
		PublicKey:      "public key1",
		StartBlock:     1000,
		LastRewardDate: exampleTime,
		FinishTime:     exampleTime,
		Status:         "Active",
		ValidatorId:    2,
	}

	example2 := structs.Node{
		Name:           "name2",
		Ip:             "127.0.0.2",
		PublicIp:       "127.0.0.2",
		Port:           1904,
		PublicKey:      "public key2",
		StartBlock:     1001,
		LastRewardDate: exampleTime,
		FinishTime:     exampleTime,
		Status:         "Pending",
		ValidatorId:    3,
	}
	exampleNodes = append(exampleNodes, example1)
	exampleNodes = append(exampleNodes, example2)

	tests := []struct {
		number     int
		name       string
		req        *http.Request
		nodes      []structs.Node
		dbResponse error
		code       int
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
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidSyntaxForNodes))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 3,
			name:   "missing parameter",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(invalidPropertyNameForNodes))),
			},
			code: http.StatusBadRequest,
		},
		{
			number: 4,
			name:   "internal server error",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForNodes))),
			},
			nodes:      exampleNodes,
			dbResponse: errors.New("internal error"),
			code:       http.StatusInternalServerError,
		},
		{
			number: 5,
			name:   "success response",
			req: &http.Request{
				Method: http.MethodPost,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(validJsonForNodes))),
			},
			nodes: exampleNodes,
			code:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := store.NewMockDataStore(mockCtrl)
			if tt.number > 3 {
				mockDB.EXPECT().SaveOrUpdateNodes(tt.req.Context(), tt.nodes).Return(tt.dbResponse)
			}
			contractor := *client.NewClientContractor(mockDB)
			connector := handler.NewClientConnector(contractor)
			res := http.HandlerFunc(connector.SaveOrUpdateNodes)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, tt.req)
			assert.True(t, rr.Code == tt.code)
		})
	}
}

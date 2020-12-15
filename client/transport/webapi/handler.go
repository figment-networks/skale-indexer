package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"go.uber.org/zap"
)

type ClientContractor interface {
	GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error)
	GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error)
	GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error)
	GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error)
	GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error)
}

// Connector is main HTTP connector for manager
type Connector struct {
	cli ClientContractor
}

// NewConnector is  Connector constructor
func NewClientConnector(cli ClientContractor) *Connector {
	return &Connector{cli}
}

func (c *Connector) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Connector) GetContractEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	params := structs.QueryParams{}
	from := req.URL.Query().Get("from")
	timeFrom, errFrom := time.Parse(structs.Layout, from)
	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(structs.Layout, to)
	if errFrom == nil && errTo == nil {
		params.TimeFrom = timeFrom
		params.TimeTo = timeTo
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
		return
	}

	typeParam := req.URL.Query().Get("type")
	idParam := req.URL.Query().Get("id")
	var err error
	if (typeParam == "" ) !=  (idParam == "") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}

	if typeParam == ""  {
		params.BoundType = typeParam
		var id uint64
		id, err = strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
		params.BoundId = append(params.BoundId, id)
	}

	res, err := c.cli.GetContractEvents(req.Context(), params)
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(newApiError(err, http.StatusNotFound))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetNodes(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	validatorIdParam := req.URL.Query().Get("validator_id")
	var validatorId uint64
	var err error
	if validatorIdParam != "" {
		validatorId, err = strconv.ParseUint(validatorIdParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}
	recentParam := req.URL.Query().Get("recent")
	recent, _ := strconv.ParseBool(recentParam)
	params := structs.QueryParams{
		Id:          id,
		ValidatorId: validatorId,
		Recent:      recent,
	}
	res, err := c.cli.GetNodes(req.Context(), params)
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(newApiError(err, http.StatusNotFound))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	validatorIdParam := req.URL.Query().Get("validator_id")
	var validatorId uint64
	var err error
	if validatorIdParam != "" {
		validatorId, err = strconv.ParseUint(validatorIdParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	from := req.URL.Query().Get("from")
	timeFrom, errFrom := time.Parse(structs.Layout, from)
	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(structs.Layout, to)
	recentParam := req.URL.Query().Get("recent")
	recent, _ := strconv.ParseBool(recentParam)
	if id == "" && validatorIdParam == "" && ((errFrom != nil || errTo != nil) || (from == "" && to == "")) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:          id,
		ValidatorId: validatorId,
		TimeFrom:    timeFrom,
		TimeTo:      timeTo,
		Recent:      recent,
	}

	res, err := c.cli.GetValidators(req.Context(), params)
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(newApiError(err, http.StatusNotFound))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetDelegations(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	validatorIdParam := req.URL.Query().Get("validator_id")
	var validatorId uint64
	var err error
	if validatorIdParam != "" {
		validatorId, err = strconv.ParseUint(validatorIdParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	from := req.URL.Query().Get("from")
	timeFrom, errFrom := time.Parse(structs.Layout, from)
	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(structs.Layout, to)
	recentParam := req.URL.Query().Get("recent")
	recent, _ := strconv.ParseBool(recentParam)

	if id == "" && validatorIdParam == "" && (errFrom != nil || errTo != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:          id,
		ValidatorId: validatorId,
		TimeFrom:    timeFrom,
		TimeTo:      timeTo,
		Recent:      recent,
	}

	res, err := c.cli.GetDelegations(req.Context(), params)
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(newApiError(err, http.StatusNotFound))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetValidatorStatistics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	validatorIdParam := req.URL.Query().Get("validator_id")
	var validatorId uint64
	var err error
	if validatorIdParam != "" {
		validatorId, err = strconv.ParseUint(validatorIdParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	statisticTypeParam := req.URL.Query().Get("statistic_type")
	if statisticTypeParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
		return
	}

	var statisticType structs.StatisticTypeVS
	if statisticTypeParam == "total_stake" {
		statisticType = structs.ValidatorStatisticsTypeTotalStake
	} else if statisticTypeParam == "active_nodes" {
		statisticType = structs.ValidatorStatisticsTypeActiveNodes
	} else if statisticTypeParam == "linked_nodes" {
		statisticType = structs.ValidatorStatisticsTypeLinkedNodes
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:              id,
		ValidatorId:     validatorId,
		StatisticTypeVS: statisticType,
	}

	res, err := c.cli.GetValidatorStatistics(req.Context(), params)
	if err != nil {
		if errors.Is(err, structs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(newApiError(err, http.StatusNotFound))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/health", c.HealthCheck)
	mux.HandleFunc("/contract-events", c.GetContractEvents)
	mux.HandleFunc("/nodes", c.GetNodes)
	mux.HandleFunc("/validators", c.GetValidators)
	mux.HandleFunc("/delegations", c.GetDelegations)
	mux.HandleFunc("/validator-statistics", c.GetValidatorStatistics)
}

type ScrapeContractor interface {
	ParseLogs(ctx context.Context, ccs map[common.Address]contract.ContractsContents, from, to big.Int) error
}

// ScrapeConnector is main HTTP connector for manager
type ScrapeConnector struct {
	l   *zap.Logger
	cli ScrapeContractor
	ccs map[common.Address]contract.ContractsContents
}

// NewConnector is  Connector constructor
func NewScrapeConnector(l *zap.Logger, sc ScrapeContractor, ccs map[common.Address]contract.ContractsContents) *ScrapeConnector {
	return &ScrapeConnector{l, sc, ccs}
}

// AttachToHandler attaches handlers to http server's mux
func (sc *ScrapeConnector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/getLogs", sc.GetLogs)
}

func (sc *ScrapeConnector) GetLogs(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	f := req.URL.Query().Get("from")
	from, ok := new(big.Int).SetString(f, 10)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"from parameters are incorrect"}`))
		return
	}

	t := req.URL.Query().Get("to")
	to, ok2 := new(big.Int).SetString(t, 10)
	if !ok2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":" to parameters are incorrect"}`))
		return
	}

	if err := sc.cli.ParseLogs(req.Context(), sc.ccs, *from, *to); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
}

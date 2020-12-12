package handler

import (
	"encoding/json"
	"errors"
	"github.com/figment-networks/skale-indexer/client/structs"
	"net/http"
	"strconv"
	"time"
)

const Layout = time.RFC3339

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
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	params := structs.QueryParams{
		Id: id,
	}
	res, err := c.cli.GetContractEvents(req.Context(), params)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
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

	params := structs.QueryParams{
		Id:          id,
		ValidatorId: validatorId,
	}
	res, err := c.cli.GetNodes(req.Context(), params)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
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
	timeFrom, errFrom := time.Parse(Layout, from)
	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(Layout, to)

	if id == "" && validatorIdParam == "" && ((errFrom != nil || errTo != nil) || (from == "" && to == "")) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:          id,
		ValidatorId: validatorId,
		TimeFrom:    timeFrom,
		TimeTo:      timeTo,
	}

	res, err := c.cli.GetValidators(req.Context(), params)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
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
	timeFrom, errFrom := time.Parse(Layout, from)
	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(Layout, to)

	if id == "" && validatorIdParam == "" && (errFrom != nil || errTo != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:          id,
		ValidatorId: validatorId,
		TimeFrom:    timeFrom,
		TimeTo:      timeTo,
	}

	res, err := c.cli.GetDelegations(req.Context(), params)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
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
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
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
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:              id,
		ValidatorId:     validatorId,
		StatisticTypeVS: statisticType,
	}

	res, err := c.cli.GetValidatorStatistics(req.Context(), params)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
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

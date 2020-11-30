package handler

import (
	"encoding/json"
	"errors"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/structs"
	"net/http"
	"strconv"
	"time"
)

const Layout = time.RFC3339

// Connector is main HTTP connector for manager
type Connector struct {
	cli client.ClientContractor
}

// NewConnector is  Connector constructor
func NewClientConnector(cli client.ClientContractor) *Connector {
	return &Connector{cli}
}

func (c *Connector) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Connector) SaveOrUpdateDelegations(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var delegations []structs.Delegation
	err := decoder.Decode(&delegations)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}
	for _, dlg := range delegations {
		err = validateDelegationRequiredFields(dlg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	err = c.cli.SaveOrUpdateDelegations(req.Context(), delegations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
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
	holderParam := req.URL.Query().Get("holder")
	var holder uint64
	if holderParam != "" {
		holder, err = strconv.ParseUint(holderParam, 10, 64)
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

	if id == "" && holderParam == "" && validatorIdParam == "" && (errFrom != nil || errTo != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:          id,
		Holder:      holder,
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

func (c *Connector) SaveOrUpdateEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var events []structs.Event
	err := decoder.Decode(&events)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}
	for _, dlg := range events {
		err = validateEventRequiredFields(dlg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	err = c.cli.SaveOrUpdateEvents(req.Context(), events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Connector) GetEvents(w http.ResponseWriter, req *http.Request) {
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
	res, err := c.cli.GetEvents(req.Context(), params)
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

func (c *Connector) SaveOrUpdateValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var validators []structs.Validator
	err := decoder.Decode(&validators)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}
	for _, vld := range validators {
		err = validateValidatorRequiredFields(vld)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	err = c.cli.SaveOrUpdateValidators(req.Context(), validators)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Connector) GetValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	address, ok := req.URL.Query()["address"]
	a := make([]structs.Address, len(address))
	for i, v := range address {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
			return
		}
		a[i] = structs.Address(val)
	}
	from := req.URL.Query().Get("from")
	timeFrom, errFrom := time.Parse(Layout, from)
	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(Layout, to)

	if id == "" && (!ok || len(address) == 0) && ((errFrom != nil || errTo != nil) || (from == "" && to == "")) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:       id,
		Address:  a,
		TimeFrom: timeFrom,
		TimeTo:   timeTo,
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

func (c *Connector) SaveOrUpdateNodes(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	decoder := json.NewDecoder(req.Body)
	var nodes []structs.Node
	err := decoder.Decode(&nodes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}
	for _, n := range nodes {
		err = validateNodeRequiredFields(n)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	err = c.cli.SaveOrUpdateNodes(req.Context(), nodes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
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

func (c *Connector) GetDelegationStatistics(w http.ResponseWriter, req *http.Request) {
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
	if statisticTypeParam == "" || !(statisticTypeParam == "states" || statisticTypeParam == "next-epoch") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	statisticType := structs.StatesStatisticsType
	if statisticTypeParam != "states" {
		statisticType = structs.NextEpochStatisticsType
	}

	params := structs.QueryParams{
		Id:            id,
		ValidatorId:   validatorId,
		StatisticType: statisticType,
	}

	res, err := c.cli.GetDelegationStatistics(req.Context(), params)
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

	mux.HandleFunc("/save-or-update-delegations", c.SaveOrUpdateDelegations)
	mux.HandleFunc("/delegations", c.GetDelegations)

	mux.HandleFunc("/save-or-update-events", c.SaveOrUpdateEvents)
	mux.HandleFunc("/events", c.GetEvents)

	mux.HandleFunc("/save-or-update-validators", c.SaveOrUpdateValidators)
	mux.HandleFunc("/validators", c.GetValidators)

	mux.HandleFunc("/save-or-update-nodes", c.SaveOrUpdateNodes)
	mux.HandleFunc("/nodes", c.GetNodes)

	mux.HandleFunc("/delegation-statistics", c.GetDelegationStatistics)
}

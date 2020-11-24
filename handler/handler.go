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

const Layout = "2006-01-02T15:04:05.000Z"

// Connector is main HTTP connector for manager
type Connector struct {
	cli client.ClientContractor
}

// NewConnector is  Connector constructor
func NewClientConnector(cli client.ClientContractor) *Connector {
	return &Connector{cli}
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

func (c *Connector) GetDelegationById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	res, err := c.cli.GetDelegationById(req.Context(), id)
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

func (c *Connector) GetDelegationsByHolder(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	holder := req.URL.Query().Get("holder")
	if holder == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	res, err := c.cli.GetDelegationsByHolder(req.Context(), holder)
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

func (c *Connector) GetDelegationsByValidatorId(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	validatorIdParam := req.URL.Query().Get("validator_id")
	validatorId, err := strconv.ParseUint(validatorIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}

	res, err := c.cli.GetDelegationsByValidatorId(req.Context(), validatorId)
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

func (c *Connector) GetEventById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	res, err := c.cli.GetEventById(req.Context(), id)
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

func (c *Connector) GetEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	res, err := c.cli.GetAllEvents(req.Context())
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
	/*
		validatorAddress := req.URL.Query().Get("address")
		if validatorAddress == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
			return
		}*/

	id := req.URL.Query().Get("id")
	from := req.URL.Query().Get("from")
	timeFrom, errFrom := time.Parse(Layout, from)

	to := req.URL.Query().Get("to")
	timeTo, errTo := time.Parse(Layout, to)

	if id == "" && (errFrom != nil || errTo != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(ErrMissingParameter, http.StatusBadRequest))
		return
	}

	params := structs.QueryParams{
		Id:       id,
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

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/save-or-update-delegations", c.SaveOrUpdateDelegations)
	mux.HandleFunc("/get-delegation", c.GetDelegationById)
	mux.HandleFunc("/get-delegations-by-holder", c.GetDelegationsByHolder)
	mux.HandleFunc("/get-delegations-by-validator-id", c.GetDelegationsByValidatorId)

	mux.HandleFunc("/save-or-update-events", c.SaveOrUpdateEvents)
	mux.HandleFunc("/get-event-by-id", c.GetEventById)
	mux.HandleFunc("/get-events", c.GetEvents)

	mux.HandleFunc("/save-or-update-validators", c.SaveOrUpdateValidators)
	mux.HandleFunc("/validators", c.GetValidators)
}

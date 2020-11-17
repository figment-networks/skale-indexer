package handler

import (
	"../client"
	"../structs"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrMissingParameter = errors.New("missing parameter")
)

// Connector is main HTTP connector for manager
type Connector struct {
	cli client.ClientContractor
}

// NewConnector is  Connector constructor
func NewClientConnector(cli client.ClientContractor) *Connector {
	return &Connector{cli}
}

func (c *Connector) SaveOrUpdateDelegations(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var delegations []structs.Delegation
	err := decoder.Decode(&delegations)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	for _, dlg := range delegations {
		err = validateDelegationRequiredFields(dlg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	err = c.cli.SaveOrUpdateDelegations(req.Context(), delegations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Connector) GetDelegationById(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := c.cli.GetDelegationById(req.Context(), &id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetDelegationsByHolder(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	holder := req.URL.Query().Get("holder")
	if holder == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrMissingParameter.Error()))
		return
	}

	res, err := c.cli.GetDelegationsByHolder(req.Context(), &holder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetDelegationsByValidatorId(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	validatorIdParam := req.URL.Query().Get("validator-id")
	validatorId, err := strconv.ParseUint(validatorIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := c.cli.GetDelegationsByValidatorId(req.Context(), &validatorId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) SaveOrUpdateValidators(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var validators []structs.Validator
	err := decoder.Decode(&validators)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	for _, vld := range validators {
		err = validateValidatorRequiredFields(vld)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	err = c.cli.SaveOrUpdateValidators(req.Context(), validators)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Connector) GetValidatorById(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := c.cli.GetValidatorById(req.Context(), &id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetValidatorsByValidatorAddress(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	validatorAddress := req.URL.Query().Get("validator_address")
	if validatorAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrMissingParameter.Error()))
		return
	}

	res, err := c.cli.GetValidatorsByValidatorAddress(req.Context(), &validatorAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

func (c *Connector) GetValidatorsByRequestedAddress(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestedAddress := req.URL.Query().Get("requested_address")
	if requestedAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrMissingParameter.Error()))
		return
	}

	res, err := c.cli.GetValidatorsByRequestedAddress(req.Context(), &requestedAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(res)
}

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/save-delegations", c.SaveOrUpdateDelegations)
	mux.HandleFunc("/get-delegation", c.GetDelegationById)
	mux.HandleFunc("/get-delegations-by-holder", c.GetDelegationsByHolder)
	mux.HandleFunc("/get-delegations-by-validator-id", c.GetDelegationsByValidatorId)
	mux.HandleFunc("/save-validators", c.SaveOrUpdateValidators)
	mux.HandleFunc("/get-validator", c.GetValidatorById)
	mux.HandleFunc("/get-validators-by-validator-address", c.GetValidatorsByValidatorAddress)
	mux.HandleFunc("/get-validators-by-requested-address", c.GetValidatorsByRequestedAddress)
}

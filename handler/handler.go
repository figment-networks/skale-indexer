package handler

import (
	"../client"
	"../structs"
	"../types"
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

func (c *Connector) SaveOrUpdateDelegation(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var delegation structs.Delegation
	err := decoder.Decode(&delegation)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = validateRequiredFields(delegation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = c.cli.SaveOrUpdateDelegation(req.Context(), delegation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
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
		err = validateRequiredFields(dlg)
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

	idParam := req.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	idType := types.ID(id)
	res, err := c.cli.GetDelegationById(req.Context(), &idType)
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

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/save-delegation", c.SaveOrUpdateDelegation)
	mux.HandleFunc("/save-delegations", c.SaveOrUpdateDelegations)
	mux.HandleFunc("/get-delegation", c.GetDelegationById)
	mux.HandleFunc("/get-delegations-by-holder", c.GetDelegationsByHolder)
	mux.HandleFunc("/get-delegations-by-validator-id", c.GetDelegationsByValidatorId)
}

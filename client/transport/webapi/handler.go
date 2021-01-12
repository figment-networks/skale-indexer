package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"go.uber.org/zap"
)

// ClientContractor - method signatures for Connector
type ClientContractor interface {
	GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error)
	GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error)
	GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error)
	GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetValidatorStatisticsTimeline(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error)
}

// Connector is main HTTP connector for manager
type Connector struct {
	cli ClientContractor
}

// NewClientConnector is connector constructor
func NewClientConnector(cli ClientContractor) *Connector {
	return &Connector{cli}
}

/**
 * Health check endpoint
 *
 * Method: any method
 * Success 200
**/
func (c *Connector) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

/**
 * Contract events endpoint
 *
 * Method: GET, POST
 * Params:
 *   see EventParams
 *   required:
 *	   @from: the inclusive beginning of the time range for event time
 *     @to: the inclusive ending of the time range for event time
 *   optional:
 *     @type: Event type (required when id is provided)
 *     @id: Bound id (required when type is provided)
 *
 * Error:
 *     http code: 400, 405, 500
 *     response: see apiError struct
 *
 * Success
 *     http code: 200
 *     response: see ContractEvent struct
**/
func (c *Connector) GetContractEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := EventParams{}
	switch req.Method {
	case http.MethodGet:
		m, err := pathParams(strings.Replace(req.URL.Path, "/event/", "", -1))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
		from := req.URL.Query().Get("from")
		to := req.URL.Query().Get("to")
		typeParam := req.URL.Query().Get("type")
		idParam := req.URL.Query().Get("id")
		if m != nil {
			if f, ok := m["from"]; ok {
				from = f
			}
			if t, ok := m["to"]; ok {
				to = t
			}
			if t, ok := m["type"]; ok {
				typeParam = t
			}
			if i, ok := m["id"]; ok {
				idParam = i
			}
		}

		timeFrom, errFrom := time.Parse(structs.Layout, from)
		timeTo, errTo := time.Parse(structs.Layout, to)
		if errFrom == nil && errTo == nil {
			params.TimeFrom = timeFrom
			params.TimeTo = timeTo
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
			return
		}

		if (typeParam == "") != (idParam == "") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
			return
		}

		if typeParam != "" {
			params.Type = typeParam
			params.Id, err = strconv.ParseUint(idParam, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("id parameter given in wrong format"), http.StatusBadRequest))
				return
			}
		}
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	res, err := c.cli.GetContractEvents(req.Context(), structs.EventParams{
		Id:       params.Id,
		Type:     params.Type,
		TimeFrom: params.TimeFrom,
		TimeTo:   params.TimeTo,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var ceva []ContractEvent
	for _, r := range res {
		ceva = append(ceva, ContractEvent{
			ID:              r.ID,
			ContractName:    r.ContractName,
			ContractAddress: r.ContractAddress,
			EventName:       r.EventName,
			BlockHeight:     r.BlockHeight,
			Time:            r.Time,
			TransactionHash: r.TransactionHash,
			Params:          r.Params,
			Removed:         r.Removed,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(ceva)
}

/**
 * Nodes endpoint
 *
 * Method: GET, POST
 * Params:
 *   see NodeParams
 *   optional:
 *     @id: the index of node in SKALE deployed smart contract
 *     @validator_id: the index of validator in SKALE deployed smart contract
 *
 * Error:
 *     http code: 400, 405, 500
 *     response: see apiError struct
 *
 * Success:
 *     http code: 200
 *     response: see Node struct
**/
func (c *Connector) GetNode(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := NodeParams{}
	switch req.Method {
	case http.MethodGet:
		m, err := pathParams(strings.Replace(req.URL.Path, "/node/", "", -1))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}

		params.NodeID = req.URL.Query().Get("id")
		params.ValidatorID = req.URL.Query().Get("validator_id")
		params.Status = req.URL.Query().Get("status")

		if m != nil {
			if nodeId, ok := m["id"]; ok {
				params.NodeID = nodeId
			}
			if validatorId, ok := m["validator_id"]; ok {
				params.ValidatorID = validatorId
			}
			if status, ok := m["status"]; ok {
				params.Status = status
			}
		}
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}
	nParams := structs.NodeParams{
		NodeID:      params.NodeID,
		ValidatorID: params.ValidatorID,
	}
	if params.Status != "" {
		var ok bool
		var s structs.NodeStatus
		if s, ok = structs.GetTypeForNode(params.Status); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("node type is wrong"), http.StatusBadRequest))
			return
		}
		nParams.Status = s.String()
	}

	res, err := c.cli.GetNodes(req.Context(), nParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var nodes []Node
	for _, n := range res {
		nodes = append(nodes, Node{
			NodeID:         n.NodeID,
			Name:           n.Name,
			IP:             n.IP.String(),
			PublicIP:       n.PublicIP.String(),
			Port:           n.Port,
			StartBlock:     n.StartBlock,
			NextRewardDate: n.NextRewardDate,
			LastRewardDate: n.LastRewardDate,
			FinishTime:     n.FinishTime,
			ValidatorID:    n.ValidatorID,
			Status:         n.Status.String(),
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(nodes)
}

/**
 * Validators endpoint
 *
 * Method: GET, POST
 * Params:
 *   see ValidatorParams
 *   optional:
 *     @id: the index of validator in SKALE deployed smart contract
 *     @from: the inclusive beginning of the time range for registration time
 *     @to: the inclusive ending of the time range for registration time
 *
 * Error:
 *     http code: 400, 405, 500
 *     response: see apiError struct
 *
 * Success:
 *     http code: 200
 *     response: see Validator struct
**/
func (c *Connector) GetValidator(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	m, err := pathParams(strings.Replace(req.URL.Path, "/validator/", "", -1))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}

	params := ValidatorParams{}
	switch req.Method {
	case http.MethodGet:
		var timeFrom, timeTo string

		params.ValidatorID = req.URL.Query().Get("id")
		timeFrom = req.URL.Query().Get("from")
		timeTo = req.URL.Query().Get("to")

		if m != nil {
			if id, ok := m["id"]; ok {
				params.ValidatorID = id
			}
			timeFrom, _ = m["from"]
			timeTo, _ = m["to"]
		}
		var errFrom, errTo error
		if !(timeFrom == "" && timeTo == "") {
			params.TimeFrom, errFrom = time.Parse(structs.Layout, timeFrom)
			params.TimeTo, errTo = time.Parse(structs.Layout, timeTo)
		}

		if errFrom != nil || errTo != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("error parsing time format (from/to) parameters"), http.StatusBadRequest))
			return
		}
	case http.MethodOptions:
		// TODO(lukanus): add options preflight headers
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	res, err := c.cli.GetValidators(req.Context(), structs.ValidatorParams{
		ValidatorID: params.ValidatorID,
		TimeFrom:    params.TimeFrom,
		TimeTo:      params.TimeTo,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(errors.New("error during server query"), http.StatusInternalServerError))
		return
	}

	var vlds []Validator
	for _, vld := range res {
		vlds = append(vlds, Validator{
			ValidatorID:             vld.ValidatorID,
			Name:                    vld.Name,
			ValidatorAddress:        vld.ValidatorAddress,
			RequestedAddress:        vld.RequestedAddress,
			Description:             vld.Description,
			FeeRate:                 vld.FeeRate,
			RegistrationTime:        vld.RegistrationTime,
			MinimumDelegationAmount: vld.MinimumDelegationAmount,
			AcceptNewRequests:       vld.AcceptNewRequests,
			Authorized:              vld.Authorized,
			ActiveNodes:             vld.ActiveNodes,
			LinkedNodes:             vld.LinkedNodes,
			Staked:                  vld.Staked.String(),
			Pending:                 vld.Pending,
			Rewards:                 vld.Rewards,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(vlds); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

/**
 * Validator statistics endpoint
 *
 * Method: GET, POST
 * Params:
 *   see ValidatorStatisticsParams
 *   optional:
 *     @id: the index of validator in SKALE deployed smart contract
 *     @type: statistics type
 *     @timeline: returns whether the latest or statistics changes timeline
 *
 * Error:
 *     http code: 400, 405, 500
 *     response: see apiError struct
 *
 * Success:
 *     http code: 200
 *     response: see ValidatorStatistic struct
**/
func (c *Connector) GetValidatorStatistics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	m, err := pathParams(strings.Replace(req.URL.Path, "/validator/statistics/", "", -1))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(err, http.StatusBadRequest))
		return
	}

	params := ValidatorStatisticsParams{}
	switch req.Method {
	case http.MethodGet:
		params.ValidatorID = req.URL.Query().Get("id")
		params.Type = req.URL.Query().Get("type")
		params.Timeline = (req.URL.Query().Get("timeline") != "")
		if m != nil {
			if id, ok := m["id"]; ok {
				params.ValidatorID = id
			}
			if typ, ok := m["type"]; ok {
				params.Type = typ
			}
			if _, ok := m["timeline"]; ok {
				params.Timeline = true
			}
		}

		if params.Timeline && params.ValidatorID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("validator id must be provided for timeline"), http.StatusBadRequest))
			return
		}
	case http.MethodOptions:
		// TODO(lukanus): add options preflight headers
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	vParams := structs.ValidatorStatisticsParams{
		ValidatorID: params.ValidatorID,
	}

	if params.Type != "" || params.Timeline {
		var ok bool
		if vParams.Type, ok = structs.GetTypeForValidatorStatistics(params.Type); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("statistic type is wrong"), http.StatusBadRequest))
			return
		}
	}

	var res []structs.ValidatorStatistics
	if params.Timeline {
		res, err = c.cli.GetValidatorStatisticsTimeline(req.Context(), vParams)
	} else {
		res, err = c.cli.GetValidatorStatistics(req.Context(), vParams)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var vlds []ValidatorStatistic
	for _, v := range res {
		vlds = append(vlds, ValidatorStatistic{
			Type:        v.Type.String(),
			ValidatorID: v.ValidatorID,
			BlockHeight: v.BlockHeight,
			Amount:      v.Amount.String(),
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(vlds); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

/**
 * Accounts endpoint
 *
 * Method: GET, POST
 * Params:
 *   see AccountParams
 *   optional:
 *     @type: account type
 *     @address: account address
 *
 * Error:
 *     http code: 400, 405, 500
 *     response: see apiError struct
 *
 * Success:
 *     http code: 200
 *     response: see Account struct
**/
func (c *Connector) GetAccount(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	params := structs.AccountParams{}
	switch req.Method {
	case http.MethodGet:
		m, err := pathParams(strings.Replace(req.URL.Path, "/account/", "", -1))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
		params.Type = req.URL.Query().Get("type")
		params.Address = req.URL.Query().Get("address")

		if m != nil {
			if t, ok := m["type"]; ok {
				params.Type = t
			}
			if address, ok := m["address"]; ok {
				params.Address = address
			}
		}
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	res, err := c.cli.GetAccounts(req.Context(), structs.AccountParams{
		Address: params.Address,
		Type:    params.Type,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var accs []Account
	for _, a := range res {
		accs = append(accs, Account{
			Address: a.Address,
			Type:    string(a.Type),
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(accs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

/**
 * Delegations endpoint
 *
 * Method: GET, POST
 * Params:
 *   see DelegationParams
 *   optional:
 *     @id: the index of delegation in SKALE deployed smart contract
 *     @validator_id: the index of validator in SKALE deployed smart contract
 *     @from: the inclusive beginning of the time range for delegation created time
 *     @to: the inclusive ending of the time range for delegation created time
 *     @timeline: returns whether the latest or delegation changes timeline
 *
 * Error:
 *     http code: 400, 405, 500
 *     response: see apiError struct
 *
 * Success:
 *     http code: 200
 *     response: see Delegation struct
**/
func (c *Connector) GetDelegation(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	params := DelegationParams{}
	switch req.Method {
	case http.MethodGet:
		m, err := pathParams(strings.Replace(req.URL.Path, "/delegation/", "", -1))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
		from := req.URL.Query().Get("from")
		to := req.URL.Query().Get("to")
		vID := req.URL.Query().Get("validator_id")
		dID := req.URL.Query().Get("id")
		params.Timeline = (req.URL.Query().Get("timeline") != "")
		if m != nil {
			if f, ok := m["from"]; ok {
				from = f
			}
			if t, ok := m["to"]; ok {
				to = t
			}
			if v, ok := m["validator_id"]; ok {
				vID = v
			}
			if d, ok := m["id"]; ok {
				dID = d
			}
			if _, ok := m["timeline"]; ok {
				params.Timeline = true
			}
		}
		params.ValidatorID = vID
		params.DelegationID = dID

		var errFrom, errTo error
		if from != "" && to != "" {
			params.TimeFrom, errFrom = time.Parse(structs.Layout, from)
			params.TimeTo, errTo = time.Parse(structs.Layout, to)
		}
		if errFrom != nil || errTo != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
			return
		}

	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	dParams := structs.DelegationParams{
		ValidatorID:  params.ValidatorID,
		DelegationID: params.DelegationID,
		TimeFrom:     params.TimeFrom,
		TimeTo:       params.TimeTo,
	}

	var (
		res []structs.Delegation
		err error
	)
	if params.Timeline {
		res, err = c.cli.GetDelegationTimeline(req.Context(), dParams)
	} else {
		res, err = c.cli.GetDelegations(req.Context(), dParams)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var dlgs []Delegation
	for _, dlg := range res {
		dlgs = append(dlgs, Delegation{
			DelegationID:    dlg.DelegationID,
			TransactionHash: dlg.TransactionHash,
			Holder:          dlg.Holder,
			ValidatorID:     dlg.ValidatorID,
			BlockHeight:     dlg.BlockHeight,
			Amount:          dlg.Amount.String(),
			Period:          dlg.DelegationPeriod,
			Started:         dlg.Started,
			Created:         dlg.Created,
			Finished:        dlg.Finished,
			Info:            dlg.Info,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(dlgs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/health", c.HealthCheck)
	mux.HandleFunc("/event/", c.GetContractEvents)
	mux.HandleFunc("/node/", c.GetNode)
	mux.HandleFunc("/validator/", c.GetValidator)
	mux.HandleFunc("/validator/statistics/", c.GetValidatorStatistics)
	mux.HandleFunc("/delegation/", c.GetDelegation)
	mux.HandleFunc("/account/", c.GetAccount)
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

// NewScrapeConnector is  Connector constructor
func NewScrapeConnector(l *zap.Logger, sc ScrapeContractor, ccs map[common.Address]contract.ContractsContents) *ScrapeConnector {
	return &ScrapeConnector{l, sc, ccs}
}

// AttachToHandler attaches handlers to http server's mux
func (sc *ScrapeConnector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/getLogs", sc.GetLogs)
}

/*
 * Gets logs from node endpoint
 */
func (sc *ScrapeConnector) GetLogs(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write(newApiError(, http.StatusMethodNotAllowed))
		w.Write([]byte(`{"error":"from parameters are incorrect"}`))
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

func pathParams(path string) (map[string]string, error) {
	p := strings.Split(path, "/")
	var p2 []string
	for _, k := range p {
		if k != "" {
			p2 = append(p2, k)
		}
	}

	switch len(p2) {
	case 0:
		return nil, nil
	case 1:
		return map[string]string{"id": p2[0]}, nil
	default:
		if len(p2)%2 == 1 {
			return nil, errors.New("path has to be in key/value pair format")
		}
		a := map[string]string{}
		for k, v := range p2 {
			if k%2 == 0 {
				a[v] = p2[k+1]
			}
		}
		return a, nil
	}

}

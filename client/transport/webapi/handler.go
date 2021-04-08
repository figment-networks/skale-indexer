package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
)

//go:generate swagger generate spec --scan-models -o swagger.json

// ClientContractor - method signatures for Connector
type ClientContractor interface {
	GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error)
	GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error)
	GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetValidatorStatisticsTimeline(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error)

	GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error)
	GetSystemEvents(ctx context.Context, params structs.SystemEventParams) (systemEvents []structs.SystemEvent, err error)
}

// Connector is main HTTP connector for manager
type Connector struct {
	cli ClientContractor
}

// NewClientConnector is connector constructor
func NewClientConnector(cli ClientContractor) *Connector {
	return &Connector{cli}
}

func (c *Connector) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetContractEvents
func (c *Connector) GetContractEvents(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	params := EventParams{}
	switch req.Method {
	case http.MethodGet:
		allowCORSHeaders(w)
		m := map[string]string{}
		var err error
		if strings.Index(req.URL.Path[1:], "/") > 0 {
			m, err = pathParams(strings.Replace(req.URL.Path, "/events/", "", -1), "id")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(err, http.StatusBadRequest))
				return
			}
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

		limit := req.URL.Query().Get("limit")
		if limit != "" {
			if params.Limit, err = strconv.ParseUint(limit, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'limit' parameter"), http.StatusBadRequest))
				return
			}
			offset := req.URL.Query().Get("offset")
			if params.Offset, err = strconv.ParseUint(offset, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'offset' parameter"), http.StatusBadRequest))
				return
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
			params.ID, err = strconv.ParseUint(idParam, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("id parameter given in wrong format"), http.StatusBadRequest))
				return
			}
		}
	case http.MethodPost:
		allowCORSHeaders(w)
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(newApiError(structs.ErrMissingParameter, http.StatusInternalServerError))
			return
		}
	case http.MethodOptions:
		allowCORSHeaders(w)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	res, err := c.cli.GetContractEvents(req.Context(), structs.EventParams{
		Id:       params.ID,
		Type:     params.Type,
		TimeFrom: params.TimeFrom,
		TimeTo:   params.TimeTo,
		Offset:   params.Offset,
		Limit:    params.Limit,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	ceva := ContractEvents{}
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

func (c *Connector) GetNode(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := NodeParams{}

	allowCORSHeaders(w)
	switch req.Method {
	case http.MethodGet:
		m := map[string]string{}
		var err error
		if strings.Index(req.URL.Path[1:], "/") > 0 {
			m, err = pathParams(strings.Replace(req.URL.Path, "/node/", "", -1), "id")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(err, http.StatusBadRequest))
				return
			}
		}
		params.NodeID = req.URL.Query().Get("id")
		params.ValidatorID = req.URL.Query().Get("validator_id")
		params.Status = req.URL.Query().Get("status")

		limit := req.URL.Query().Get("limit")
		if limit != "" {
			if params.Limit, err = strconv.ParseUint(limit, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'limit' parameter"), http.StatusBadRequest))
				return
			}
			offset := req.URL.Query().Get("offset")
			if params.Offset, err = strconv.ParseUint(offset, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'offset' parameter"), http.StatusBadRequest))
				return
			}
		}

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
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}
	nParams := structs.NodeParams{
		NodeID:      params.NodeID,
		ValidatorID: params.ValidatorID,
		Offset:      params.Offset,
		Limit:       params.Limit,
	}
	if params.Status != "" {
		var ok bool
		if _, ok = structs.GetTypeForNode(params.Status); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("node type is wrong"), http.StatusBadRequest))
			return
		}
		nParams.Status = params.Status
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
			Address:        n.Address,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(nodes)
}

func (c *Connector) GetValidator(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	allowCORSHeaders(w)

	m := map[string]string{}
	var err error
	if strings.Index(req.URL.Path[1:], "/") > 0 {
		m, err = pathParams(strings.Replace(req.URL.Path, "/validators/", "", -1), "id")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	params := ValidatorParams{}
	switch req.Method {
	case http.MethodGet:
		params.ValidatorID = req.URL.Query().Get("id")
		timeFrom := req.URL.Query().Get("from")
		timeTo := req.URL.Query().Get("to")

		authorized := req.URL.Query().Get("authorized")
		if authorized != "" {
			auth, err := strconv.ParseUint(authorized, 10, 64)
			if err != nil || auth > 2 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing authorized"), http.StatusBadRequest))
				return
			}
			params.Authorized = uint8(auth)
		}

		limit := req.URL.Query().Get("limit")
		if limit != "" {
			if params.Limit, err = strconv.ParseUint(limit, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'limit' parameter"), http.StatusBadRequest))
				return
			}
			offset := req.URL.Query().Get("offset")
			if params.Offset, err = strconv.ParseUint(offset, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'offset' parameter"), http.StatusBadRequest))
				return
			}
		}
		if m != nil {
			if id, ok := m["id"]; ok {
				params.ValidatorID = id
			}
			if f, ok := m["from"]; ok {
				timeFrom = f
			}
			if t, ok := m["to"]; ok {
				timeTo = t
			}
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
		return
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
		Authorized:  structs.ThreeState(params.Authorized),
		Offset:      params.Offset,
		Limit:       params.Limit,
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
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(vlds); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

func (c *Connector) GetValidatorStatistics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	allowCORSHeaders(w)

	m := map[string]string{}
	var err error
	if strings.Index(req.URL.Path[1:], "/") > 0 {
		m, err = pathParams(strings.Replace(req.URL.Path, "/validators/statistics/", "", -1), "id")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	params := ValidatorStatisticsParams{}
	switch req.Method {
	case http.MethodGet:
		from := req.URL.Query().Get("from")
		to := req.URL.Query().Get("to")
		params.ValidatorID = req.URL.Query().Get("id")
		params.Type = req.URL.Query().Get("type")
		params.Timeline = (req.URL.Query().Get("timeline") != "")
		if m != nil {
			if f, ok := m["from"]; ok {
				from = f
			}
			if t, ok := m["to"]; ok {
				to = t
			}
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

		limit := req.URL.Query().Get("limit")
		if limit != "" {
			if params.Limit, err = strconv.ParseUint(limit, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'limit' parameter"), http.StatusBadRequest))
				return
			}
			offset := req.URL.Query().Get("offset")
			if params.Offset, err = strconv.ParseUint(offset, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'offset' parameter"), http.StatusBadRequest))
				return
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

		if params.Timeline && params.ValidatorID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("validator id must be provided for timeline"), http.StatusBadRequest))
			return
		}
	case http.MethodOptions:
		return
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
		TimeFrom:    params.TimeFrom,
		TimeTo:      params.TimeTo,
		Limit:       params.Limit,
		Offset:      params.Offset,
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
		vld := ValidatorStatistic{
			Type:        v.Type.String(),
			ValidatorID: v.ValidatorID,
			BlockHeight: v.BlockHeight,
			BlockTime:   v.Time,
			Amount:      v.Amount.String(),
		}
		if v.Type == structs.ValidatorStatisticsTypeValidatorAddress || v.Type == structs.ValidatorStatisticsTypeRequestedAddress {
			vld.Amount = common.ToHex(v.Amount.Bytes())
		}

		vlds = append(vlds, vld)
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(vlds); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

func (c *Connector) GetAccount(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	allowCORSHeaders(w)

	params := structs.AccountParams{}
	switch req.Method {
	case http.MethodGet:
		m := map[string]string{}
		var err error
		if strings.Index(req.URL.Path[1:], "/") > 0 {
			m, err = pathParams(strings.Replace(req.URL.Path, "/accounts/", "", -1), "id")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(err, http.StatusBadRequest))
				return
			}
		}
		params.Type = req.URL.Query().Get("type")
		params.Address = req.URL.Query().Get("address")

		limit := req.URL.Query().Get("limit")
		if limit != "" {
			if params.Limit, err = strconv.ParseUint(limit, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'limit' parameter"), http.StatusBadRequest))
				return
			}
			offset := req.URL.Query().Get("offset")
			if params.Offset, err = strconv.ParseUint(offset, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'offset' parameter"), http.StatusBadRequest))
				return
			}
		}

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
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	res, err := c.cli.GetAccounts(req.Context(), structs.AccountParams{
		Address: params.Address,
		Type:    params.Type,
		Limit:   params.Limit,
		Offset:  params.Offset,
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

func (c *Connector) GetDelegation(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	allowCORSHeaders(w)

	params := DelegationParams{}
	switch req.Method {
	case http.MethodGet:
		m := map[string]string{}
		var err error
		if strings.Index(req.URL.Path[1:], "/") > 0 {
			m, err = pathParams(strings.Replace(req.URL.Path, "/delegations/", "", -1), "id")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(err, http.StatusBadRequest))
				return
			}
		}

		state := req.URL.Query().Get("state")
		if state != "" {
			state = strings.TrimPrefix(state, "[")
			state = strings.TrimSuffix(state, "]")
			params.State = strings.Split(state, ",")
		}

		limit := req.URL.Query().Get("limit")
		if limit != "" {
			if params.Limit, err = strconv.ParseUint(limit, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'limit' parameter"), http.StatusBadRequest))
				return
			}
			offset := req.URL.Query().Get("offset")
			if params.Offset, err = strconv.ParseUint(offset, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'offset' parameter"), http.StatusBadRequest))
				return
			}
		}

		from := req.URL.Query().Get("from")
		to := req.URL.Query().Get("to")
		vID := req.URL.Query().Get("validator_id")
		dID := req.URL.Query().Get("id")
		holder := req.URL.Query().Get("holder")
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
			if h, ok := m["holder"]; ok {
				holder = h
			}
			if _, ok := m["timeline"]; ok {
				params.Timeline = true
			}
		}
		params.ValidatorID = vID
		params.DelegationID = dID
		params.Holder = holder

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
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	// DelegationState
	var dss []structs.DelegationState
	for _, st := range params.State {
		ds := structs.DelegationStateFromString(st)
		if ds == structs.DelegationStateUNKNOWN {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(errors.New("wrong delegation state literal"), http.StatusBadRequest))
			return
		}
		dss = append(dss, ds)
	}

	dParams := structs.DelegationParams{
		ValidatorID:  params.ValidatorID,
		DelegationID: params.DelegationID,
		State:        dss,
		Holder:       params.Holder,
		TimeFrom:     params.TimeFrom,
		TimeTo:       params.TimeTo,
		Offset:       params.Offset,
		Limit:        params.Limit,
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
			State:           dlg.State.String(),
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(dlgs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

func (c *Connector) GetSystemEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	allowCORSHeaders(w)

	m := map[string]string{}
	var err error
	if strings.Index(req.URL.Path[1:], "/") > 0 {
		m, err = pathParams(strings.Replace(req.URL.Path, "/system_events/", "", -1), "address")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}
	params := SystemEventParams{}
	switch req.Method {
	case http.MethodGet:
		params.Address = req.URL.Query().Get("address")
		params.Kind = req.URL.Query().Get("kind")
		after := req.URL.Query().Get("after")

		if m != nil {
			if ad, ok := m["address"]; ok {
				params.Address = ad
			}
			if k, ok := m["kind"]; ok {
				params.Kind = k
			}
			if a, ok := m["after"]; ok {
				after = a
			}
		}

		if after != "" {
			if params.After, err = strconv.ParseUint(after, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(newApiError(errors.New("error parsing 'after' parameter"), http.StatusBadRequest))
				return
			}
		}
	case http.MethodOptions:
		return
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

	res, err := c.cli.GetSystemEvents(req.Context(), structs.SystemEventParams{
		After:      params.After,
		Kind:       params.Kind,
		Address:    params.Address,
		SenderID:   params.SenderID,
		ReceiverID: params.ReceiverID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(errors.New("error during server query"), http.StatusInternalServerError))
		return
	}

	var sEvts []SystemEvent
	for _, evt := range res {
		sevt := structs.SysEvtTypes[evt.Kind]
		sEvts = append(sEvts, SystemEvent{
			Height:      evt.Height,
			Time:        evt.Time,
			Kind:        sevt,
			Sender:      evt.Sender,
			Recipient:   evt.Recipient,
			SenderID:    evt.SenderID.Uint64(),
			RecipientID: evt.RecipientID.Uint64(),
			Data: SystemEventData{
				After:  evt.After,
				Before: evt.Before,
			},
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	if err := enc.Encode(sEvts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
	}
}

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/health", c.HealthCheck)

	// swagger:operation GET /events Event getContractEvents
	//
	// Contract events endpoint
	//
	// This endpoint returns events that comes from  SKALE ethereum contracts
	//
	// ---
	// Produces:
	// - application/json
	// Schemes:
	// - http
	//
	// Parameters:
	//   - in: query
	//     name: from
	//     x-go-type:
	//       import:
	//         package: "time"
	//     required: true
	//     type: string
	//     description: the inclusive beginning of the time range for event time
	//   - in: query
	//     name: to
	//     x-go-type:
	//       import:
	//         package: "time"
	//     required: true
	//     type: string
	//     description: the inclusive ending of the time range for event time
	//   - in: query
	//     name: type
	//     type: string
	//     required: false
	//     description: event type
	//     example: validator
	//   - in: query
	//     name: id
	//     type: string
	//     required: false
	//     description: bound id
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/ContractEvents"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"

	// swagger:operation POST /events Event getContractEvents
	//
	// Contract events endpoint
	//
	// This endpoint returns events that comes from  SKALE ethereum contracts
	//
	// ---
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes:
	// - http
	//
	// Parameters:
	// - name: EventParams
	//   schema:
	//     "$ref": "#/definitions/EventParams"
	//   in: body
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/ContractEvents"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	mux.HandleFunc("/events/", c.GetContractEvents)
	mux.HandleFunc("/events", c.GetContractEvents)

	// swagger:operation GET /node Nodes getNodes
	//
	// Node returning endpoint
	//
	// This endpoint returns node information
	//
	// ---
	// Produces:
	// - application/json
	// Schemes:
	// - http
	//
	// Parameters:
	//   - in: query
	//     name: id
	//     type: string
	//     required: false
	//     description: the index of node in SKALE deployed smart contract
	//   - in: query
	//     name: validator_id
	//     type: string
	//     required: false
	//     description: the index of validator in SKALE deployed smart contract
	//   - in: query
	//     name: status
	//     type: string
	//     required: false
	//     description: node status
	//     example: Active
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Nodes"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"

	// swagger:operation POST /node Nodes getNodes
	//
	// Node returning endpoint
	//
	// This endpoint returns node information
	//
	// ---
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes:
	// - http
	//
	// Parameters:
	// - name: NodeParams
	//   schema:
	//     "$ref": "#/definitions/NodeParams"
	//   in: body
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Nodes"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	mux.HandleFunc("/nodes/", c.GetNode)
	mux.HandleFunc("/nodes", c.GetNode)

	// swagger:operation GET /validators Validator getValidators
	//
	// Validators returning endpoint
	//
	// This endpoint returns validator information
	//
	// ---
	// Produces:
	// - application/json
	// Schemes:
	// - http
	//
	// Parameters:
	//   - in: query
	//     name: id
	//     type: string
	//     required: false
	//     description: the index of validator in SKALE deployed smart contract
	//   - in: query
	//     name: from
	//     x-go-type:
	//       import:
	//         package: "time"
	//     type: string
	//     required: false
	//     description: the inclusive beginning of the time range for registration time
	//   - in: query
	//     name: to
	//     x-go-type:
	//       import:
	//         package: "time"
	//     type: string
	//     required: false
	//     description: the inclusive ending of the time range for registration time
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Validators"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"

	// swagger:operation POST /validators Validator getValidators
	//
	// Validator returning endpoint
	//
	// This endpoint returns Validator information
	//
	// ---
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes:
	// - http
	//
	// Parameters:
	// - name: ValidatorParams
	//   schema:
	//     "$ref": "#/definitions/ValidatorParams"
	//   in: body
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Validators"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	mux.HandleFunc("/validators/", c.GetValidator)
	mux.HandleFunc("/validators", c.GetValidator)

	// swagger:operation GET /validators/statistics ValidatorStatistics getValidatorStatistics
	//
	// Validator statistics returning endpoint
	//
	// This endpoint returns validator statistics information
	//
	// ---
	// Produces:
	// - application/json
	// Schemes:
	// - http
	//
	// Parameters:
	//   - in: query
	//     name: id
	//     type: string
	//     required: false
	//     description: the index of validator in SKALE deployed smart contract
	//   - in: query
	//     name: type
	//     type: string
	//     required: false
	//     description: statistics type
	//     example: TOTAL_STAKE
	//   - in: query
	//     name: timeline
	//     type: boolean
	//     required: false
	//     description: returns whether the latest or statistics changes timeline
	//   - in: query
	//     name: from
	//     x-go-type:
	//       import:
	//         package: "time"
	//     type: string
	//     required: true
	//     description: the inclusive beginning of the time range for block time
	//   - in: query
	//     name: to
	//     x-go-type:
	//       import:
	//         package: "time"
	//     type: string
	//     required: true
	//     description: the inclusive ending of the time range for block time
	//
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/ValidatorStatistics"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"

	// swagger:operation POST /validators/statistics ValidatorStatistics getValidatorStatistics
	//
	// Validator statistics returning endpoint
	//
	// This endpoint returns validator statistics information
	//
	// ---
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes:
	// - http
	//
	// Parameters:
	// - name: ValidatorStatisticsParams
	//   schema:
	//     "$ref": "#/definitions/ValidatorStatisticsParams"
	//   in: body
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/ValidatorStatistics"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	mux.HandleFunc("/validators/statistics/", c.GetValidatorStatistics)
	mux.HandleFunc("/validators/statistics", c.GetValidatorStatistics)

	// swagger:operation GET /delegations Delegations getDelegations
	//
	// Delegations returning endpoint
	//
	// This endpoint returns delegation information
	//
	// ---
	// Produces:
	// - application/json
	// Schemes:
	// - http
	//
	// Parameters:
	//   - in: query
	//     name: id
	//     type: string
	//     required: false
	//     description: the index of delegation in SKALE deployed smart contract
	//   - in: query
	//     name: validator_id
	//     type: string
	//     required: false
	//     description: the index of validator in SKALE deployed smart contract
	//   - in: query
	//     name: holder
	//     type: string
	//     description: holder address
	//     required: false
	//   - in: query
	//     name: timeline
	//     type: boolean
	//     required: false
	//     description: returns whether the latest or delegation changes timeline
	//   - in: query
	//     name: from
	//     x-go-type:
	//       import:
	//         package: "time"
	//     type: string
	//     required: false
	//     description: the inclusive beginning of the time range for delegation created time
	//   - in: query
	//     name: to
	//     x-go-type:
	//       import:
	//         package: "time"
	//     type: string
	//     required: false
	//     description: the inclusive ending of the time range for delegation created time
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Delegations"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"

	// swagger:operation POST /delegations Delegations getDelegations
	//
	// Delegation returning endpoint
	//
	// This endpoint returns delegation information
	//
	// ---
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes:
	// - http
	//
	// Parameters:
	// - name: DelegationParams
	//   schema:
	//     "$ref": "#/definitions/DelegationParams"
	//   in: body
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Delegations"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	mux.HandleFunc("/delegations/", c.GetDelegation)
	mux.HandleFunc("/delegations", c.GetDelegation)

	// swagger:operation GET /accounts Account getAccounts
	//
	// Accounts returning endpoint
	//
	// This endpoint returns account information
	//
	// ---
	// Produces:
	// - application/json
	// Schemes:
	// - http
	//
	// Parameters:
	//   - in: query
	//     name: type
	//     type: string
	//     description: account type
	//     required: false
	//     example: delegator
	//   - in: query
	//     name: address
	//     type: string
	//     description:  account address
	//     required: false
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Accounts"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"

	// swagger:operation POST /accounts Account getAccounts
	//
	// Account returning endpoint
	//
	// This endpoint returns account information
	//
	// ---
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes:
	// - http
	//
	// Parameters:
	// - name: DelegationParams
	//   schema:
	//     "$ref": "#/definitions/AccountParams"
	//   in: body
	//
	// Responses:
	//   default:
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '200':
	//     schema:
	//       "$ref": "#/definitions/Accounts"
	//   '400':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	//   '500':
	//     schema:
	//       "$ref": "#/definitions/ApiError"
	mux.HandleFunc("/accounts/", c.GetAccount)
	mux.HandleFunc("/accounts", c.GetAccount)

	mux.HandleFunc("/system_events/", c.GetSystemEvents)
	mux.HandleFunc("/system_events", c.GetSystemEvents)
}

func pathParams(path, key string) (map[string]string, error) {
	p := strings.Split(path, "/")
	p2 := []string{}
	for _, k := range p {
		if k != "" {
			p2 = append(p2, k)
		}
	}

	switch len(p2) {
	case 0:
		return nil, nil
	case 1:
		return map[string]string{key: p2[0]}, nil
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

func allowCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
}

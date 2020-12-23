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
	GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error)
	GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error)
	GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error)
	GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error)
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

	params := structs.EventParams{}
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
			w.Write(newApiError(err, http.StatusBadRequest))
			return
		}
	}

	res, err := c.cli.GetContractEvents(req.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var ceva []ContractEventAPI
	for _, r := range res {
		ceva = append(ceva, ContractEventAPI{
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

func (c *Connector) GetNodes(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	validatorId := req.URL.Query().Get("validator_id")
	params := structs.NodeParams{
		ValidatorId: validatorId,
	}
	res, err := c.cli.GetNodes(req.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var nodes []NodeAPI
	for _, n := range res {
		nodes = append(nodes, NodeAPI{
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
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(nodes)
}

func (c *Connector) GetValidators(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	params := structs.ValidatorParams{
		ValidatorId: req.URL.Query().Get("validator_id"),
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

	var vlds []ValidatorAPI
	for _, vld := range res {
		vlds = append(vlds, ValidatorAPI{
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
			Active:                  vld.Active,
			ActiveNodes:             vld.ActiveNodes,
			LinkedNodes:             vld.LinkedNodes,
			Staked:                  vld.Staked,
			Pending:                 vld.Pending,
			Rewards:                 vld.Rewards,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(vlds)
}

func (c *Connector) GetDelegations(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	timeFrom, errFrom := time.Parse(structs.Layout, req.URL.Query().Get("from"))
	timeTo, errTo := time.Parse(structs.Layout, req.URL.Query().Get("to"))

	if errFrom != nil || errTo != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newApiError(structs.ErrMissingParameter, http.StatusBadRequest))
		return
	}

	res, err := c.cli.GetDelegations(req.Context(), structs.DelegationParams{
		ValidatorId:  req.URL.Query().Get("validator_id"),
		DelegationId: req.URL.Query().Get("delegation_id"),
		TimeFrom:     timeFrom,
		TimeTo:       timeTo,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var dlgs []DelegationAPI
	for _, dlg := range res {
		dlgs = append(dlgs, DelegationAPI{
			DelegationID:     dlg.DelegationID,
			Holder:           dlg.Holder,
			ValidatorID:      dlg.ValidatorID,
			BlockHeight:      dlg.BlockHeight,
			Amount:           dlg.Amount,
			DelegationPeriod: dlg.DelegationPeriod,
			Created:          dlg.Created,
			Info:             dlg.Info,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(dlgs)
}

// TODO: add unit tests
func (c *Connector) GetValidatorStatistics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}
	recentParam := req.URL.Query().Get("recent")
	recent, _ := strconv.ParseBool(recentParam)
	params := structs.ValidatorStatisticsParams{
		ValidatorId:     req.URL.Query().Get("validator_id"),
		StatisticTypeVS: req.URL.Query().Get("statistic_type"),
		Recent:          recent,
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

	var vlds []ValidatorStatisticsAPI
	for _, v := range res {
		vlds = append(vlds, ValidatorStatisticsAPI{
			StatisticType: v.StatisticType.String(),
			ValidatorId:   v.ValidatorId.Uint64(),
			BlockHeight:   v.BlockHeight,
			Amount:        v.Amount,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(vlds)
}

func (c *Connector) GetAccounts(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(newApiError(structs.ErrNotAllowedMethod, http.StatusMethodNotAllowed))
		return
	}

	params := structs.AccountParams{
		Type: req.URL.Query().Get("type"),
	}

	res, err := c.cli.GetAccounts(req.Context(), params)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	var accs []AccountAPI
	for _, a := range res {
		accs = append(accs, AccountAPI{
			Address:     a.Address,
			AccountType: a.AccountType,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	enc.Encode(accs)
}

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/health", c.HealthCheck)
	mux.HandleFunc("/contract-events", c.GetContractEvents)
	mux.HandleFunc("/nodes", c.GetNodes)
	mux.HandleFunc("/validators", c.GetValidators)
	mux.HandleFunc("/delegations", c.GetDelegations)
	mux.HandleFunc("/validator-statistics", c.GetValidatorStatistics)
	mux.HandleFunc("/accounts", c.GetAccounts)
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

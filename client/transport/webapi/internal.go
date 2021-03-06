package webapi

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ScrapeContractor interface {
	ParseLogs(ctx context.Context, taskID string, from, to big.Int) error
	GetLatestData(ctx context.Context, taskID string, latest uint64) (lastHeight uint64, isRunning bool, err error)
}

// ScrapeConnector is main HTTP connector for manager
type ScrapeConnector struct {
	l                   *zap.Logger
	cli                 ScrapeContractor
	scrapeLatestTimeout time.Duration
}

// NewScrapeConnector is  Connector constructor
func NewScrapeConnector(l *zap.Logger, sc ScrapeContractor, scrapeLatestTimeout time.Duration) *ScrapeConnector {
	return &ScrapeConnector{l, sc, scrapeLatestTimeout}
}

// AttachToHandler attaches handlers to http server's mux
func (sc *ScrapeConnector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/getLogs", sc.GetLogs)
	mux.HandleFunc("/scrape_latest", sc.GetLatest)
}

/*
 * Gets logs from node endpoint
 */
func (sc *ScrapeConnector) GetLogs(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"method not allowed"}`))
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

	if err := sc.cli.ParseLogs(req.Context(), uuid.New().String(), *from, *to); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newApiError(err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
}

type LatestDataRequest struct {
	Network string `json:"network"`
	ChainID string `json:"chain_id"`
	Version string `json:"version"`

	TaskID     string `json:"task_id"`
	LastHeight uint64 `json:"last_height"`
}

type LatestDataResponse struct {
	LastHeight uint64 `json:"last_height"`
	Error      []byte `json:"error"`
	Processing bool   `json:"processing"`
}

/*
 * Gets latest entries after certain height
 */
func (sc *ScrapeConnector) GetLatest(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	lDResp := LatestDataResponse{}
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		lDResp.Error = []byte(`{"error":"method not allowed"}`)
		enc.Encode(lDResp)
		return
	}

	dec := json.NewDecoder(req.Body)
	ldr := &LatestDataRequest{}
	if err := dec.Decode(ldr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lDResp.Error = []byte(`{"error":"error decoding LatestDataRequest format "}`)
		enc.Encode(lDResp)
		return
	}

	tCtx, cancel := context.WithTimeout(req.Context(), sc.scrapeLatestTimeout)
	defer cancel()
	// TODO(lukanus): Check the version
	lastHeight, isRunning, err := sc.cli.GetLatestData(tCtx, ldr.TaskID, ldr.LastHeight)
	if err != nil {
		lDResp.LastHeight = ldr.LastHeight
		w.WriteHeader(http.StatusInternalServerError)
		lDResp.Error = []byte(err.Error())
		enc.Encode(lDResp)
		return
	}

	if isRunning == true {
		w.WriteHeader(http.StatusOK)
		lDResp.Processing = true
		enc.Encode(lDResp)
		return
	}

	lDResp.LastHeight = lastHeight
	w.WriteHeader(http.StatusOK)
	err = enc.Encode(lDResp)
	if err != nil {
		sc.l.Error("Error encoding response  ", zap.Error(err))
	}
}

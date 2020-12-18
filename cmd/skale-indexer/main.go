package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"net/http"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/skale-indexer/api/skale"
	"github.com/figment-networks/skale-indexer/client"
	"github.com/figment-networks/skale-indexer/client/actions"
	"github.com/figment-networks/skale-indexer/client/transport/webapi"
	"github.com/figment-networks/skale-indexer/cmd/skale-indexer/config"
	"github.com/figment-networks/skale-indexer/cmd/skale-indexer/logger"
	"github.com/figment-networks/skale-indexer/scraper"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/store/postgresql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.Parse()
}

func main() {
	ctx := context.Background()


	// connect to database
	db, err := sql.Open("postgres", "postgresql://localhost:5432/skale?user=postgres&password=admin&sslmode=disable")
	if err != nil {
		return
	}

	if err := db.PingContext(ctx); err != nil {
		logger.Error(err)
		return
	}
	defer db.Close()

	pgsqlDriver := postgresql.NewDriver(ctx, db, logger.GetLogger())
	storeDB := store.New(pgsqlDriver)

	tr := eth.NewEthTransport("http://0.0.0.0:8545")
	if err := tr.Dial(ctx); err != nil { // TODO(lukanus): check if this has recovery
		return
	}
	defer tr.Close(ctx)

	cm := contract.NewManager()
	if err := cm.LoadContractsFromDir("C:\\Users\\emire\\repo\\skale-indexer\\test\\integration\\testFIles"); err != nil {
		return
	}

	caller := &skale.Caller{}
	am := actions.NewManager(caller, storeDB, tr, cm, nil)
	eAPI := scraper.NewEthereumAPI(logger.GetLogger(), tr, am)
	mux := http.NewServeMux()

	cli := client.NewClient(storeDB)
	hCli := webapi.NewClientConnector(cli)
	hCli.AttachToHandler(mux)

	ccs := cm.GetContractsByNames(am.GetImplementedContractNames())
	sCli := webapi.NewScrapeConnector(logger.GetLogger(), eAPI, ccs)
	sCli.AttachToHandler(mux)

	mux.Handle("/metrics", metrics.Handler())

	s := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: mux,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if err := s.ListenAndServe(); err != nil {
		logger.GetLogger().Error("[HTTP] failed to listen", zap.Error(err))
	}

}

func initConfig(path string) (*config.Config, error) {
	cfg := &config.Config{}
	if path != "" {
		if err := config.FromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

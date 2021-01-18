package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/metrics/prometheusmetrics"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Initialize configuration
	cfg, err := initConfig(configPath)
	if err != nil {
		log.Fatalf("error initializing config [ERR: %v]", err.Error())
	}

	if cfg.RollbarServerRoot == "" {
		cfg.RollbarServerRoot = "github.com/figment-networks/skale-indexer"
	}
	rcfg := &logger.RollbarConfig{
		AppEnv:             cfg.AppEnv,
		RollbarAccessToken: cfg.RollbarAccessToken,
		RollbarServerRoot:  cfg.RollbarServerRoot,
		Version:            config.GitSHA,
	}

	if cfg.AppEnv == "development" || cfg.AppEnv == "local" {
		logger.Init("console", "debug", []string{"stderr"}, rcfg)
	} else {
		logger.Init("json", "info", []string{"stderr"}, rcfg)
	}

	logger.Info(config.IdentityString())
	defer logger.Sync()

	// Initialize metrics
	prom := prometheusmetrics.New()
	err = metrics.AddEngine(prom)
	if err != nil {
		logger.Error(err)
	}
	err = metrics.Hotload(prom.Name())
	if err != nil {
		logger.Error(err)
	}

	// connect to database
	logger.Info("[DB] Connecting to database...")
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error(err)
		return
	}

	if err := db.PingContext(ctx); err != nil {
		logger.Error(err)
		return
	}
	logger.Info("[DB] Ping successfull...")
	defer db.Close()

	pgsqlDriver := postgresql.NewDriver(ctx, db, logger.GetLogger())
	storeDB := store.New(pgsqlDriver)

	mux := http.NewServeMux()
	cli := client.NewClient(logger.GetLogger(), storeDB)
	hCli := webapi.NewClientConnector(cli)
	hCli.AttachToHandler(mux)

	if cfg.EnableScraper {
		logger.GetLogger().Info("Indexer is in scraping mode")
		caller := &skale.Caller{}
		nodeTypeMessage := "Ethereum node is in archive mode"
		if cfg.EthereumNodeType == "recent" {
			caller.NodeType = skale.ENTRecent
			nodeTypeMessage = "Ethereum node is in recent mode"
		}
		logger.GetLogger().Info(nodeTypeMessage)

		cm := contract.NewManager()
		if err := cm.LoadContractsFromDir(cfg.SkaleABIDir); err != nil {
			logger.Fatal("Error dialing", zap.String("directory", cfg.SkaleABIDir), zap.Error(err))
			return
		}
		tr := eth.NewEthTransport(cfg.EthereumAddress)
		if err := tr.Dial(ctx); err != nil { // TODO(lukanus): check if this has recovery
			logger.Fatal("Error dialing ethereum", zap.String("ethereum_address", cfg.EthereumAddress), zap.Error(err))
			return
		}
		defer tr.Close(ctx)
		am := actions.NewManager(caller, storeDB, tr, cm, logger.GetLogger())
		eAPI := scraper.NewEthereumAPI(logger.GetLogger(), tr, am)
		ccs := cm.GetContractsByNames(am.GetImplementedContractNames())
		sCli := webapi.NewScrapeConnector(logger.GetLogger(), eAPI, ccs)
		sCli.AttachToHandler(mux)
	} else {
		logger.GetLogger().Info("Indexer is not in scraping mode")
	}

	mux.Handle("/metrics", metrics.Handler())

	s := &http.Server{
		Addr:    cfg.Address,
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

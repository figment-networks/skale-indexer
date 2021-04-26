package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
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

	"github.com/figment-networks/indexing-engine/health"
	"github.com/figment-networks/indexing-engine/health/database/postgreshealth"

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

	if cfg.EnableScraper {
		logger.GetLogger().Info("Indexer is in scraping mode")
		caller := skale.NewCaller(skale.ENTArchive, cfg.RequestsPerSecond)
		nodeTypeMessage := "Ethereum node is in archive mode"
		if cfg.EthereumNodeType == "recent" {
			caller.NodeType = skale.ENTRecent
			nodeTypeMessage = "Ethereum node is in recent mode"
		}
		logger.GetLogger().Info(nodeTypeMessage)

		cm := contract.NewManager()

		if cfg.AdditionalABI != "" {
			f, err := os.Open(cfg.AdditionalABI)
			if err != nil {
				logger.Fatal("Error opening default events file", zap.Error(err))
				return
			}

			if err = cm.AddGlobalEvents(f); err != nil {
				f.Close()
				logger.Fatal("Error loading default events ", zap.Error(err))
				return
			}
			f.Close()
		}

		if err := cm.LoadContractsFromDir(cfg.SkaleABIDir); err != nil {
			logger.Fatal("Error getting contracts", zap.String("directory", cfg.SkaleABIDir), zap.Error(err))
			return
		}
		logger.GetLogger().Info("Loaded contracts", zap.String("dir", cfg.SkaleABIDir))

		tr := eth.NewEthTransport(cfg.EthereumAddress)
		if err := tr.Dial(ctx); err != nil { // TODO(lukanus): check if this has recovery
			logger.Fatal("Error dialing ethereum", zap.String("ethereum_address", cfg.EthereumAddress), zap.Error(err))
			return
		}
		defer tr.Close(ctx)
		am := actions.NewManager(caller, storeDB, tr, cm, logger.GetLogger())
		eAPI := scraper.NewEthereumAPI(logger.GetLogger(), tr, types.Header{Number: new(big.Int).SetUint64(cfg.EthereumSmallestBlockNumber), Time: cfg.EthereumSmallestTime}, am)

		cli := client.NewClient(logger.GetLogger(),
			storeDB, eAPI,
			cm.GetContractsByNames(am.GetImplementedContractNames()),
			cfg.EthereumSmallestBlockNumber,
			cfg.MaxHeightsPerRequest)
		hCli := webapi.NewClientConnector(cli)
		hCli.AttachToHandler(mux)

		sCli := webapi.NewScrapeConnector(logger.GetLogger(), cli, cfg.ScrapeLatestTimeout)
		sCli.AttachToHandler(mux)
	} else {
		logger.GetLogger().Info("Indexer is not in scraping mode")

		cli := client.NewClient(logger.GetLogger(), storeDB, nil, nil, cfg.EthereumSmallestBlockNumber, cfg.MaxHeightsPerRequest)
		hCli := webapi.NewClientConnector(cli)
		hCli.AttachToHandler(mux)
	}

	mux.Handle("/metrics", metrics.Handler())

	dbMonitor := postgreshealth.NewPostgresMonitorWithMetrics(db, logger.GetLogger())
	monitor := &health.Monitor{}
	monitor.AddProber(ctx, dbMonitor)
	go monitor.RunChecks(ctx, cfg.HealthCheckInterval)
	monitor.AttachHttp(mux)

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
			return cfg, err
		}
	}

	if cfg.DatabaseURL != "" {
		return cfg, nil
	}

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

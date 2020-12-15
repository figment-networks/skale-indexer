package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/figment-networks/skale-indexer/cmd/skale-indexer-migration/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type flags struct {
	configPath    string
	migrationPath string
	version       uint
	verbose       bool
}

var configFlags = flags{}

func init() {
	flag.StringVar(&configFlags.configPath, "config", "", "Path to config")
	flag.BoolVar(&configFlags.verbose, "verbose", true, "Verbosity of logs during run")
	flag.UintVar(&configFlags.version, "version", 0, "Version parameter sets db changes to specified revision (up or down)")
	flag.StringVar(&configFlags.migrationPath, "path", "./migrations", "Path to migration folder")
	flag.Parse()
}

func main() {
	log.SetOutput(os.Stdout)
	// Initialize configuration
	cfg, err := initConfig(configFlags.configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing config [ERR: %+v]", err))
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	srcPath := fmt.Sprintf("file://%s", configFlags.migrationPath)

	if configFlags.verbose {
		log.Println("Using migrations from: ", configFlags.migrationPath)
	}

	if configFlags.version > 0 {
		if configFlags.verbose {
			log.Println("Migrating to version: ", configFlags.version)
		}
		err = migrateTo(srcPath, cfg.DatabaseURL, configFlags.version)
	} else {
		err = runMigrations(srcPath, cfg.DatabaseURL)
	}

	if err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
			return
		}
		if configFlags.verbose {
			log.Println("No change")
		}
	}

}

func initConfig(path string) (config.Config, error) {
	cfg := &config.Config{}

	if path != "" {
		if err := config.FromFile(path, cfg); err != nil {
			return *cfg, err
		}
	}

	if cfg.DatabaseURL != "" {
		return *cfg, nil
	}

	if err := config.FromEnv(cfg); err != nil {
		return *cfg, err
	}

	return *cfg, nil
}

func runMigrations(srcPath, dbURL string) error {
	m, err := migrate.New(srcPath, dbURL)
	if err != nil {
		return err
	}

	defer m.Close()

	return m.Up()
}

func migrateTo(srcPath, dbURL string, version uint) error {
	m, err := migrate.New(srcPath, dbURL)
	if err != nil {
		return err
	}

	defer m.Close()

	return m.Migrate(version)
}

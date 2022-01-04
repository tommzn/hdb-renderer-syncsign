package main

import (
	"context"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	core "github.com/tommzn/hdb-core"
)

func bootstrap(conf config.Config, ctx context.Context) (*core.Minion, error) {

	if conf == nil {
		conf = loadConfig()
	}
	logger := newLogger(conf, newSecretsManager())
	server := newServer(conf, logger, newFactory(conf, logger))
	return core.NewMinion(server), nil
}

// loadConfig from config file.
func loadConfig() config.Config {

	if conf, err := config.NewFileConfigSource(nil).Load(); err == nil {
		return conf
	}

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		exitOnError(err)
	}

	conf, err := configSource.Load()
	if err != nil {
		exitOnError(err)
	}
	return conf
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager) log.Logger {
	return log.NewLoggerFromConfig(conf, secretsMenager)
}

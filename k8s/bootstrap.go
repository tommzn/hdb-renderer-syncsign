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
	logger := newLogger(conf, newSecretsManager(), ctx)
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
	secretsManager := secrets.NewDockerecretsManager("/run/secrets/token")
	secrets.ExportToEnvironment([]string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}, secretsManager)
	return secretsManager
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager, ctx context.Context) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logContextValues := make(map[string]string)
	logContextValues[log.LogCtxNamespace] = "hdb-renderer-syncsign"
	logger.WithContext(log.LogContextWithValues(ctx, logContextValues))
	return logger
}

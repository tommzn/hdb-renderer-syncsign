package main

import (
	"fmt"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func logValue(v interface{}) {
	fmt.Println(v)
}

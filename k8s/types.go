package main

import (
	"context"
	"net/http"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-renderer-core"
	syncsign "github.com/tommzn/hdb-renderer-syncsign"
)

type webServer struct {
	port           string
	minifyResponse bool
	conf           config.Config
	logger         log.Logger
	diFactory      *factory
	httpServer     *http.Server
}

type factory struct {
	conf                    config.Config
	logger                  log.Logger
	ctx                     context.Context
	errorTemplate           core.Template
	responseTemplate        core.Template
	indoorClimateTemplate   core.Template
	timestampTemplate       core.Template
	indoorClimateRenderer   core.Renderer
	indoorClimateDataSource core.DataSource
	responseRenderer        map[string]core.Renderer
	displayConfig           *syncsign.DisplayConfig
}

type emptyResponse struct {
	StatusCode int `json:"code"`
}

type testResponse struct {
	StatusCode int                `json:"code"`
	Data       []testResponseData `json:"data"`
}

type testResponseData struct {
	RenderId string              `json:"renderId"`
	NodeId   string              `json:"nodeId"`
	Content  testResponseContent `json:"content"`
}
type testResponseContent struct {
	//Background string             `json:"background"`
	Items []testResponseItem `json:"items"`
}
type testResponseItem struct {
	Id string `json:"id"`
}

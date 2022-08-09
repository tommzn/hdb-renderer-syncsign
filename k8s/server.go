package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/mux"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

func newServer(conf config.Config, logger log.Logger, diFactory *factory) *webServer {
	port := conf.Get("hdb.server.port", config.AsStringPtr("8080"))
	minify := conf.GetAsBool("hdb.server.minify", config.AsBoolPtr(true))
	return &webServer{
		port:           *port,
		minifyResponse: *minify,
		conf:           conf,
		logger:         logger,
		diFactory:      diFactory,
	}
}

// Run starts a HTTP server to listen for rendering requests.
func (server *webServer) Run(ctx context.Context, waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	defer server.logger.Flush()

	router := mux.NewRouter()
	router.Use(server.jsonContentTypeMiddleware)
	router.Use(server.logMiddleware)

	router.HandleFunc("/renders/nodes/{nodeid}", server.handleNodeRequest).Methods("GET")
	router.HandleFunc("/renders/{renderid}", server.handleRenderRequest).Methods("GET")

	router.HandleFunc("/health", server.handleHealthCheckRequest).Methods("GET")
	router.HandleFunc("/metrics", server.handleMetricsRequest).Methods("GET")

	server.logger.Infof("Listen [%s]", server.port)
	server.logger.Flush()
	server.httpServer = &http.Server{Addr: ":" + server.port, Handler: router}

	endChan := make(chan error, 1)
	go func() {
		endChan <- server.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		server.stopHttpServer()
	case err := <-endChan:
		return err
	}
	return nil
}

// StopHttpServer will try to sop running HTTP server graceful. Timeout is 3s.
func (server *webServer) stopHttpServer() {
	server.logger.Info("Stopping HTTP server.")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("Unable to stop HTTP server, reason: ", err)
	}
}

// JsonContentTypeMiddleware adds JSON content-type header
func (server *webServer) jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// LogMiddleware adds a logger for all requests. Used log level if debug.
func (server *webServer) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.logger.Debugf("Method: %s, URL: %+v, Header: %+v, URI: %s", r.Method, r.URL, r.Header, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// HandleHealthCheckRequest always returns a 204 status code.
func (server *webServer) handleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// HandleRenderRequest handles requests for render id, which will be answered with a 204 status code by default.
func (server *webServer) handleRenderRequest(w http.ResponseWriter, r *http.Request) {

	defer server.logger.Flush()

	vars := mux.Vars(r)
	if renderId, ok := vars["renderid"]; ok {
		server.logger.Infof("Receive request for render id: %s", renderId)
		res := emptyResponse{
			StatusCode: http.StatusNoContent,
		}
		json.NewEncoder(w).Encode(res)
	} else {
		msg := "No render id passed."
		server.logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// HandleNodeRequest will try to render content for passed node.
func (server *webServer) handleNodeRequest(w http.ResponseWriter, r *http.Request) {

	defer server.logger.Flush()

	vars := mux.Vars(r)
	nodeId, ok := vars["nodeid"]
	if !ok {
		server.writeResponseError(w, "<nil>", errors.New("Node id missing."))
		return
	}

	if !server.diFactory.newDisplayConfig().Exists(nodeId) {
		server.writeResponseError(w, nodeId, fmt.Errorf("Render request for unknown node %s received.", nodeId))
		return
	}

	responseRenderer := server.diFactory.newResponseRenderer(nodeId)
	content, err := responseRenderer.Content()
	if err != nil {
		server.writeResponseError(w, nodeId, err)
		return
	}
	server.writeResponse(w, content)
}

// handleMetricsRequest will collect metrics from all datasources.
func (server *webServer) handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(server.diFactory.dataSourceMetrics())
}

// WriteResponse writes given content to response writer. If statusCode is not nil it's set as header.
func (server *webServer) writeResponse(w http.ResponseWriter, content string) {
	buf := &bytes.Buffer{}
	err := json.Compact(buf, []byte(content))
	if err == nil {
		minifiedContent := buf.Bytes()
		w.Write(minifiedContent)
	} else {
		server.logger.Error("Failed to minity content, reason: ", err)
		w.Write([]byte(content))
	}
}

// WriteResponseError will generate a error response and write it to given response writer.
func (server *webServer) writeResponseError(w http.ResponseWriter, nodeId string, err error) {

	// Remove new line before JSON conversion
	re := regexp.MustCompile(`\r?\n`)
	err = errors.New(re.ReplaceAllString(err.Error(), "\\n"))

	server.logger.Error(err)
	errorRenderer := server.diFactory.newErrorResponseRenderer(nodeId, err)
	errContent, _ := errorRenderer.Content()
	server.writeResponse(w, errContent)
}

package main

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
	conf       config.Config
	logger     log.Logger
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.conf = loadConfigForTest(nil)
	suite.logger = loggerForTest()
}
func (suite *ServerTestSuite) TestHealthRequest() {

	server := newServer(suite.conf, suite.logger, newFactory(suite.conf, suite.logger))
	suite.startServer(server)

	resp, err := http.Get("http://localhost:8080/health")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusNoContent, resp.StatusCode)

	suite.stopServer()
}

func (suite *ServerTestSuite) TestRenderRequest() {

	server := newServer(suite.conf, suite.logger, newFactory(suite.conf, suite.logger))
	suite.startServer(server)

	resp1, err1 := http.Get("http://localhost:8080/renders/xYx-123-yYy")
	suite.Nil(err1)
	suite.NotNil(resp1)
	suite.Equal(http.StatusOK, resp1.StatusCode)

	resp2, err2 := http.Get("http://localhost:8080/renders/")
	suite.Nil(err2)
	suite.NotNil(resp2)
	suite.Equal(http.StatusNotFound, resp2.StatusCode)

	suite.stopServer()
}

func (suite *ServerTestSuite) TestNodeRequest() {

	server := newServer(suite.conf, suite.logger, newFactory(suite.conf, suite.logger))
	suite.startServer(server)

	resp1, err1 := http.Get("http://localhost:8080/renders/nodes/InvalidNodeId")
	suite.Nil(err1)
	suite.NotNil(resp1)
	suite.Equal(http.StatusOK, resp1.StatusCode)

	nodeId := "Display01"
	resp2, err2 := http.Get("http://localhost:8080/renders/nodes/" + nodeId)
	suite.Nil(err2)
	suite.NotNil(resp2)
	resBody2 := suite.readBody(resp2)
	//logValue(string(resBody2))

	suite.Equal(http.StatusOK, resp2.StatusCode)
	resData := suite.asTestResponse(resBody2)
	suite.Len(resData.Data, 1)
	suite.Equal(nodeId, resData.Data[0].NodeId)
	suite.Len(resData.Data[0].Content.Items, 8)
	suite.stopServer()
}

func (suite *ServerTestSuite) startServer(server *webServer) {
	suite.wg = &sync.WaitGroup{}
	suite.ctx, suite.cancelFunc = context.WithCancel(context.Background())
	go func() {
		suite.wg.Add(1)
		suite.Nil(server.Run(suite.ctx, suite.wg))
	}()
	time.Sleep(1 * time.Second)
}

func (suite *ServerTestSuite) stopServer() {

	waitChan := make(chan bool, 1)
	go func() {
		suite.wg.Wait()
		waitChan <- true
	}()

	suite.cancelFunc()
	select {
	case <-time.After(1 * time.Second):
		suite.T().Error("Server stop timeput!")
	case ok := <-waitChan:
		suite.True(ok)
	}
}

func (suite *ServerTestSuite) readBody(res *http.Response) []byte {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	suite.Nil(err)
	return body
}

func (suite *ServerTestSuite) asTestResponse(body []byte) testResponse {
	content := testResponse{}
	suite.Nil(json.Unmarshal(body, &content))
	return content
}

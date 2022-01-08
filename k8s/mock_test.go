package main

import (
	"context"
	"github.com/stretchr/testify/suite"
	hdbcore "github.com/tommzn/hdb-core"
	core "github.com/tommzn/hdb-renderer-core"
	"testing"
	"time"
)

type DataSourceMockTestSuite struct {
	suite.Suite
}

func TestDataSourceMockTestSuite(t *testing.T) {
	suite.Run(t, new(DataSourceMockTestSuite))
}

func (suite *DataSourceMockTestSuite) TestPublishRandomEvents() {

	displays := []string{"Display1", "Display3", "Display7"}
	mock := newDataSourceMock(displays, loggerForTest())
	mock.(*dataSourceMock).publishInterval = 1 * time.Second
	eventChan := mock.Observe(nil)

	latestEvent1, err := mock.Latest(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.NotNil(err)
	suite.Nil(latestEvent1)

	allEvents1, err := mock.All(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.NotNil(err)
	suite.Nil(allEvents1)

	canelFunc, endChan := runMessagePublishing(mock)

	time.Sleep(3 * time.Second)
	suite.True(len(eventChan) >= 2)

	latestEvent2, err := mock.Latest(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.Nil(err)
	suite.NotNil(latestEvent2)

	allEvents2, err := mock.All(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.Nil(err)
	suite.NotNil(allEvents2)
	suite.True(len(allEvents2) >= 2)

	canelFunc()
	suite.True(endAsExpected(endChan, 200*time.Millisecond))
}

func (suite *DataSourceMockTestSuite) TestStackLimit() {

	displays := []string{"Display1", "Display3", "Display7"}
	mock := newDataSourceMock(displays, loggerForTest())
	mock.(*dataSourceMock).publishInterval = 10 * time.Millisecond
	filter := []hdbcore.DataSource{hdbcore.DATASOURCE_INDOORCLIMATE}
	eventChan := mock.Observe(&filter)

	canelFunc, endChan := runMessagePublishing(mock)
	time.Sleep(2 * time.Second)

	suite.Len(eventChan, 100)

	latestEvent, err := mock.Latest(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.Nil(err)
	suite.NotNil(latestEvent)

	allEvents, err := mock.All(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.Nil(err)
	suite.NotNil(allEvents)
	suite.Len(allEvents, 100)

	canelFunc()
	suite.True(endAsExpected(endChan, 200*time.Millisecond))
}

func (suite *DataSourceMockTestSuite) TestInitiMessages() {

	displays := []string{"Display1", "Display3", "Display7"}
	mock := newDataSourceMock(displays, loggerForTest())

	mock.(*dataSourceMock).initMessages()

	allEvents, err := mock.All(hdbcore.DATASOURCE_INDOORCLIMATE)
	suite.Nil(err)
	suite.NotNil(allEvents)
	suite.Len(allEvents, 9)
}

func runMessagePublishing(mock core.DataSource) (context.CancelFunc, chan bool) {
	ctx, canelFunc := context.WithCancel(context.Background())
	endChan := make(chan bool, 2)
	go func() {
		mock.(*dataSourceMock).Run(ctx)
		endChan <- true
	}()
	return canelFunc, endChan
}

func endAsExpected(endChan chan bool, timeout time.Duration) bool {
	select {
	case ok := <-endChan:
		return ok
	case <-time.After(timeout):
		return false
	}
}

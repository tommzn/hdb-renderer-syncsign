package syncsign

import (
	"context"
	"github.com/stretchr/testify/suite"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"testing"
	"time"
)

type IndoorClimateTestSuite struct {
	suite.Suite
}

func TestIndoorClimateTestSuite(t *testing.T) {
	suite.Run(t, new(IndoorClimateTestSuite))
}

func (suite *IndoorClimateTestSuite) TestGenerateContent() {

	renderer := indoorClimateRendererForTest("fixtures/testconfig02.yml")

	content, err := renderer.Content()
	suite.Nil(err)
	assertTemplateHash(suite.Assert(), content, "72ed7031fc3eb4e707cd78bdeda6e0acb2039ce6")

	size := renderer.Size()
	suite.Equal(200, size.Height)
	suite.Equal(400, size.Width)
}

func (suite *IndoorClimateTestSuite) TestGenerateContentWithError() {

	renderer := indoorClimateRendererWithDataSourceErrorForTest("fixtures/testconfig02.yml")

	suite.Len(renderer.(*indoorClimateRenderer).roomClimate, 0)
	renderer.(*indoorClimateRenderer).initIndoorClimateData()
	suite.Len(renderer.(*indoorClimateRenderer).roomClimate, 0)

	content, err := renderer.Content()
	suite.Nil(err)
	suite.Equal("", content)

	renderer2 := indoorClimateRendererWithTemplateErrorForTest("fixtures/testconfig02.yml")
	content2, err2 := renderer2.Content()
	suite.NotNil(err2)
	suite.Equal("", content2)
}

func (suite *IndoorClimateTestSuite) TestDataSourceObserving() {

	renderer := indoorClimateRendererForTest("fixtures/testconfig02.yml")

	content, err := renderer.Content()
	suite.Nil(err)
	assertTemplateHash(suite.Assert(), content, "72ed7031fc3eb4e707cd78bdeda6e0acb2039ce6")

	ctx, cancelFunc := context.WithCancel(context.Background())
	go renderer.ObserveDataSource(ctx)
	time.Sleep(1 * time.Second)

	// Increate temperature which have to change generated content
	newTemperature := "26.7"
	message := &events.IndoorClimate{
		Timestamp: timestamppb.New(time.Now()),
		DeviceId:  "Device2",
		Type:      events.MeasurementType_TEMPERATURE,
		Value:     newTemperature,
	}
	renderer.(*indoorClimateRenderer).datasource.(*datasourceMock).writeToMessageChannel(message)

	time.Sleep(1 * time.Second)
	content2, err2 := renderer.Content()
	suite.Nil(err2)
	suite.True(strings.Contains(content2, newTemperature))
	assertTemplateHash(suite.Assert(), content2, "0817ef3063f0e63752784f657401fbaa8995bf69")

	cancelFunc()
}

func (suite *IndoorClimateTestSuite) TestStopDataSourceObserving() {

	renderer := indoorClimateRendererForTest("fixtures/testconfig02.yml")
	ctx, cancelFunc := context.WithCancel(context.Background())

	endChan := make(chan bool, 1)
	go func() {
		renderer.ObserveDataSource(ctx)
		endChan <- true
	}()
	time.Sleep(100 * time.Millisecond)

	cancelFunc()
	select {
	case ok := <-endChan:
		suite.True(ok)
	case <-time.After(200 * time.Millisecond):
		suite.T().Error("DataSource observing doesn't end as expected!")
	}
}

func (suite *IndoorClimateTestSuite) TestStopDataSourceObservingOnClosedChannel() {

	renderer := indoorClimateRendererForTest("fixtures/testconfig02.yml")
	ctx, _ := context.WithCancel(context.Background())

	endChan := make(chan bool, 1)
	go func() {
		renderer.ObserveDataSource(ctx)
		endChan <- true
	}()
	time.Sleep(100 * time.Millisecond)

	close(renderer.(*indoorClimateRenderer).datasource.(*datasourceMock).eventChan)
	select {
	case ok := <-endChan:
		suite.True(ok)
	case <-time.After(200 * time.Millisecond):
		suite.T().Error("DataSource observing doesn't end as expected!")
	}
}

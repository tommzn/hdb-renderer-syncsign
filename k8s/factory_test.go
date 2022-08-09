package main

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FactoryTestSuite struct {
	suite.Suite
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}

func (suite *FactoryTestSuite) TestCreate() {

	diFactory := newFactory(loadConfigForTest(nil), loggerForTest(), context.Background())

	suite.NotNil(diFactory.newResponseRendererTemplate())
	suite.NotNil(diFactory.newErrorRendererTemplate())
	suite.NotNil(diFactory.newIndoorClimateTemplate())

	suite.NotNil(diFactory.newErrorRenderer(errors.New("Error occured!")))
	suite.NotNil(diFactory.newErrorResponseRenderer("Node01", errors.New("Error occured!")))
	suite.NotNil(diFactory.newResponseRenderer("Node01"))
	suite.NotNil(diFactory.newIndoorClimateRenderer())

	suite.NotNil(diFactory.newDataSource())

	displayConfig := diFactory.newDisplayConfig()
	suite.NotNil(displayConfig)
	suite.Len(displayConfig.All(), 3)
}

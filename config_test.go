package syncsign

import (
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestDisplayCondig() {

	displayConfig := NewDisplayConfig(loadConfigForTest(config.AsStringPtr("fixtures/testconfig03.yml")))
	suite.NotNil(displayConfig)

	suite.True(displayConfig.Exists("Display1"))
	suite.False(displayConfig.Exists("Display2"))

	suite.Len(displayConfig.All(), 2)
}

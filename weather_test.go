package syncsign

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type WeatherTestSuite struct {
	suite.Suite
}

func TestWeatherTestSuite(t *testing.T) {
	suite.Run(t, new(WeatherTestSuite))
}

func (suite *WeatherTestSuite) TestGenerateContent() {

	renderer := weatherRendererForTest("fixtures/testconfig06.yml")

	content, err := renderer.Content()
	suite.Nil(err)
	assertTemplateHash(suite.Assert(), content, "b186bb59e848ffb7145132348cdfcd75a05a0ee9")
}

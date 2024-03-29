package syncsign

import (
	"github.com/stretchr/testify/suite"
	"testing"

	events "github.com/tommzn/hdb-events-go"
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
	assertTemplateHash(suite.Assert(), content, "6f8c9bfe2b5ecec87dd998bb1aa150bd08eadc56")
}

func (suite *WeatherTestSuite) TestFormatWindSpeed() {

	fixtures := weatherDataForTest()
	weatherData := fixtures[0].(*events.WeatherData)

	suite.Equal("45", formatWindSpeed(weatherData.Current))

	weatherData.Current.WindSpeed = 7
	weatherData.Current.WindGust = 32
	suite.Equal("7/32", formatWindSpeed(weatherData.Current))
}

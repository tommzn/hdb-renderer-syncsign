package syncsign

import (
	"regexp"
	"strings"
	"time"

	"crypto/sha1"
	"encoding/hex"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	core "github.com/tommzn/hdb-renderer-core"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/assert"
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

/*
func weatherRendererForTest(configFile string) *WeatherRenderer {
	datasource := newDataSourceMock(false, false, fixturesForWeatherRenderer())
	conf := loadConfigForTest(config.AsStringPtr(configFile))
	currentWeatherTemplate := templateQithFileForTest("templates/weather_current.json")
	forecastTemplate := templateQithFileForTest("templates/weather_forecast.json")
	return NewWeatherRenderer(conf, loggerForTest(), currentWeatherTemplate, forecastTemplate, datasource)
}
*/
func billingReportRendererForTest(configFile string) *BillingReportRenderer {
	datasource := newDataSourceMock(false, false, fixturesForBillingReportRenderer())
	conf := loadConfigForTest(config.AsStringPtr(configFile))
	return NewBillingReportRenderer(conf, loggerForTest(), templateQithFileForTest("templates/billingreport.json"), datasource)
}

func indoorClimateRendererForTest(configFile string) *IndoorClimateRenderer {
	datasource := newDataSourceMock(false, false, indoorClimateDataForTest())
	return NewIndoorClimateRenderer(loadConfigForTest(config.AsStringPtr(configFile)), loggerForTest(), templateForTest(), datasource)
}

func indoorClimateRendererWithDataSourceErrorForTest(configFile string) *IndoorClimateRenderer {
	datasource := newDataSourceMock(true, true, indoorClimateDataForTest())
	return NewIndoorClimateRenderer(loadConfigForTest(config.AsStringPtr(configFile)), loggerForTest(), templateForTest(), datasource)
}

func indoorClimateRendererWithTemplateErrorForTest(configFile string) *IndoorClimateRenderer {
	datasource := newDataSourceMock(false, false, indoorClimateDataForTest())
	return NewIndoorClimateRenderer(loadConfigForTest(config.AsStringPtr(configFile)), loggerForTest(), failingTemplateForTest(), datasource)
}

func templateForTest() core.Template {
	return core.NewFileTemplate("templates/indoorclimate.json")
}

func templateQithFileForTest(filename string) core.Template {
	return core.NewFileTemplate(filename)
}

func failingTemplateForTest() core.Template {
	return newFailingTemplate()
}

func indoorClimateDataForTest() map[hdbcore.DataSource][]proto.Message {

	messages := []proto.Message{
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device2",
			Type:      events.MeasurementType_BATTERY,
			Value:     "23",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device1",
			Type:      events.MeasurementType_TEMPERATURE,
			Value:     "23.5",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device1",
			Type:      events.MeasurementType_HUMIDITY,
			Value:     "57",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device2",
			Type:      events.MeasurementType_TEMPERATURE,
			Value:     "17.1",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device2",
			Type:      events.MeasurementType_HUMIDITY,
			Value:     "65",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device1",
			Type:      events.MeasurementType_BATTERY,
			Value:     "97",
		},
	}
	events := make(map[hdbcore.DataSource][]proto.Message)
	events[hdbcore.DATASOURCE_INDOORCLIMATE] = messages
	return events
}

func fixturesForBillingReportRenderer() map[hdbcore.DataSource][]proto.Message {
	events := make(map[hdbcore.DataSource][]proto.Message)
	events[hdbcore.DATASOURCE_BILLINGREPORT] = billingReportForTest()
	events[hdbcore.DATASOURCE_EXCHANGERATE] = exchangeRateForTest()
	return events
}

func exchangeRateForTest() []proto.Message {
	return []proto.Message{
		&events.ExchangeRates{
			Rates: []*events.ExchangeRate{
				&events.ExchangeRate{
					FromCurrency: "USD",
					ToCurrency:   "EUR",
					Rate:         0.8345,
					Timestamp:    timestamppb.New(time.Now()),
				},
			},
		},
	}
}

func billingReportForTest() []proto.Message {
	billingAmount := make(map[string]float64)
	taxAmount := make(map[string]float64)
	billingAmount["xxx"] = 5.14
	billingAmount["zzz"] = 12.53
	taxAmount["xxx"] = 0.87
	taxAmount["zzz"] = 2.15
	return []proto.Message{
		&events.BillingReport{
			BillingPeriod: "Jan 2022",
			BillingAmount: billingAmount,
			TaxAmount:     taxAmount,
		},
	}
}

func fixturesForWeatherRenderer() map[hdbcore.DataSource][]proto.Message {
	events := make(map[hdbcore.DataSource][]proto.Message)
	events[hdbcore.DATASOURCE_WEATHER] = weatherDataForTest()
	return events
}

func weatherDataForTest() []proto.Message {
	return []proto.Message{
		&events.WeatherData{
			Location: &events.Location{
				Longitude: 1.0,
				Latitude:  1.0,
			},
			Units: "celsius",
			Current: &events.CurrentWeather{
				Timestamp:   timestamppb.New(time.Now()),
				Temperature: 21.7,
				WindSpeed:   45.7,
				Weather: &events.WeatherDetails{
					ConditionId: 1000,
					Group:       "sunny",
					Description: "Sunny",
					Icon:        "01d",
				},
			},
			Forecast: []*events.ForecastWeather{
				&events.ForecastWeather{
					Timestamp: timestamppb.New(time.Now().Add(1 * 24 * time.Hour)),
					Temperatures: &events.ForecastTemperatures{
						Morning: 10.1,
						Day:     17.5,
						Evening: 16.4,
						Night:   12.4,
						DayMin:  14.5,
						DayMax:  21.4,
					},
					WindSpeed: 20.5,
					Weather: &events.WeatherDetails{
						ConditionId: 1000,
						Group:       "sunny",
						Description: "Sunny",
						Icon:        "01d",
					},
				},
				&events.ForecastWeather{
					Timestamp: timestamppb.New(time.Now().Add(2 * 24 * time.Hour)),
					Temperatures: &events.ForecastTemperatures{
						Morning: 10.1,
						Day:     17.5,
						Evening: 16.4,
						Night:   12.4,
						DayMin:  14.5,
						DayMax:  21.4,
					},
					WindSpeed: 20.5,
					Weather: &events.WeatherDetails{
						ConditionId: 1000,
						Group:       "sunny",
						Description: "Sunny",
						Icon:        "01d",
					},
				},
				&events.ForecastWeather{
					Timestamp: timestamppb.New(time.Now().Add(3 * 24 * time.Hour)),
					Temperatures: &events.ForecastTemperatures{
						Morning: 10.1,
						Day:     17.5,
						Evening: 16.4,
						Night:   12.4,
						DayMin:  14.5,
						DayMax:  21.4,
					},
					WindSpeed: 20.5,
					Weather: &events.WeatherDetails{
						ConditionId: 1000,
						Group:       "sunny",
						Description: "Sunny",
						Icon:        "01d",
					},
				},
			},
		},
	}
}

func assertTemplateHash(assert *assert.Assertions, template string, expectedHash string) {
	hash := sha1.New()
	hash.Write([]byte(template))
	assert.Equal(expectedHash, hex.EncodeToString(hash.Sum(nil)))
}

func replaceUUID(content, newId string) string {
	expression := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")
	matches := expression.FindStringSubmatch(content)
	for _, match := range matches {
		content = strings.ReplaceAll(content, match, newId)
	}
	return content
}

func replaceTimeStamp(content, newTimeStamp string) string {
	expression := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Z]{3}")
	matches := expression.FindStringSubmatch(content)
	for _, match := range matches {
		content = strings.ReplaceAll(content, match, newTimeStamp)
	}
	return content
}

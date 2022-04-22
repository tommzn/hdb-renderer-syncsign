package syncsign

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	core "github.com/tommzn/hdb-renderer-core"
)

// NewWeatherRenderer returns a renderer which generates items for current weather and forcast.
func NewWeatherRenderer(conf config.Config, logger log.Logger, currentWeatherTemplate core.Template, forecastTemplate core.Template, datasource core.DataSource) *WeatherRenderer {

	anchor := anchorFromConfig(conf, "hdb.weather.anchor")
	currentWeatherSize := sizeFromConfig(conf, "hdb.weather.current.size")
	forecastWeatherSize := sizeFromConfig(conf, "hdb.weather.forecast.size")
	forecastLimit := conf.GetAsInt("hdb.weather.forecast.limit", config.AsIntPtr(6))
	return &WeatherRenderer{
		currentWeatherTemplate: currentWeatherTemplate,
		forecastTemplate:       forecastTemplate,
		currentWeatherSize:     currentWeatherSize,
		forecastWeatherSize:    forecastWeatherSize,
		forecastLimit:          *forecastLimit,
		anchor:                 anchor,
		logger:                 logger,
		datasource:             datasource,
		weatherIconMap:         newIconMap(),
	}
}

// Content generates items for weather data.
func (renderer *WeatherRenderer) Content() (string, error) {

	if renderer.weatherData == nil {
		if err := renderer.fetchEvents(); err != nil {
			renderer.logger.Errorf("Unable to get weather data, reason: %s", err)
			return "", err
		}
	}

	content, err := renderer.currentWeatherTemplate.RenderWith(renderer.currentWeatherData())
	if err != nil {
		return content, err
	}

	forecastData := renderer.forecastWeatherData()
	for _, forecast := range forecastData {
		forecastContent, err := renderer.forecastTemplate.RenderWith(forecast)
		if err != nil {
			return content, err
		}
		content += forecastContent
	}
	return content, nil
}

// FetchEvents will retrieve latest weather data.
func (renderer *WeatherRenderer) fetchEvents() error {

	weather, err := renderer.datasource.Latest(hdbcore.DATASOURCE_WEATHER)
	if err == nil {
		renderer.processEvent(weather)
	}
	return err
}

// ObserveDataSource will listen for new billing reports and exchange rate events, if report and display currency differs.
func (renderer *WeatherRenderer) ObserveDataSource(ctx context.Context) {

	defer renderer.logger.Flush()

	filter := []hdbcore.DataSource{hdbcore.DATASOURCE_WEATHER}
	renderer.dataSourceChan = renderer.datasource.Observe(&filter)
	for {
		select {
		case message, ok := <-renderer.dataSourceChan:
			if !ok {
				renderer.logger.Error("Error at reading datasource channel. Stop observing!")
				return
			}
			renderer.processEvent(message)
		case <-ctx.Done():
			renderer.logger.Info("Camceled, stop observing.")
			return
		}
	}
}

// ProcessEvent will store latest billing report and exchange rates for comtemt remdering.
func (renderer *WeatherRenderer) processEvent(message proto.Message) {

	if weatherData, ok := message.(*events.WeatherData); ok {
		renderer.logger.Debug("Receive new weather data")
		renderer.weatherData = weatherData
	}
}

func (renderer *WeatherRenderer) currentWeatherData() weatherData {
	return weatherData{
		Anchor:        renderer.anchor,
		WeatherIcon:   renderer.weatherIconMap.toWeatherIcon(renderer.weatherData.Current.Weather.Icon),
		Temperature:   fmt.Sprintf("%.1f", renderer.weatherData.Current.Temperature),
		WindSpeed:     fmt.Sprintf("%d", int(renderer.weatherData.Current.WindSpeed)),
		WindDirection: degreesToDirection(renderer.weatherData.Current.WindDirection),
		Day:           renderer.weatherData.Current.Timestamp.AsTime().Format("Monday"),
		DisplayIndex:  0,
	}
}

func (renderer *WeatherRenderer) forecastWeatherData() []weatherData {

	displayIndex := 0
	anchor := renderer.anchor
	anchor.Y += renderer.currentWeatherSize.Height
	forecasts := []weatherData{}
	for _, forecast := range renderer.weatherData.Forecast {
		forecasts = append(forecasts, weatherData{
			Anchor:           anchor,
			WeatherIcon:      renderer.weatherIconMap.toWeatherIcon(forecast.Weather.Icon),
			Temperature:      fmt.Sprintf("%.1f", forecast.Temperatures.Day),
			NightTemperature: fmt.Sprintf("%.1f", forecast.Temperatures.Night),
			WindSpeed:        fmt.Sprintf("%d", int(forecast.WindSpeed)),
			Day:              forecast.Timestamp.AsTime().Format("Monday"),
			DisplayIndex:     displayIndex,
		})
		displayIndex++
		anchor.X += renderer.forecastWeatherSize.Width
	}
	if len(forecasts) > renderer.forecastLimit {
		return forecasts[:renderer.forecastLimit]
	}
	return forecasts
}

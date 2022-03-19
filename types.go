package syncsign

import (
	"github.com/golang/protobuf/proto"
	log "github.com/tommzn/go-log"
	events "github.com/tommzn/hdb-events-go"
	core "github.com/tommzn/hdb-renderer-core"
)

type IndoorClimateRenderer struct {
	originAnchor   core.Point
	size           core.Size
	spacing        core.Spacing
	datasource     core.DataSource
	template       core.Template
	logger         log.Logger
	roomClimate    map[string]indoorCliemate
	roomCfg        roomConfig
	dataSourceChan <-chan proto.Message
	timestapMgr    core.TimestampManager
}

type indoorCliemate struct {
	DisplayIndex     string
	Temperature      string
	Humidity         string
	BatteryIcon      batteryLevelIcon
	BatteryIconColor textColor
	RoomName         string
	Anchor           core.Point
}

type roomConfig struct {
	rooms     map[string]room
	deviceMap map[string]string
}

type room struct {
	Id, Name, DisplayIndex string
}

type DisplayConfig struct {
	displays map[string]struct{}
}

type TimestampRenderer struct {
	template core.Template
}

type textColor string

const (
	COLOR_WHITE textColor = "WHITE"
	COLOR_BLACK textColor = "BLACK"
	COLOR_RED   textColor = "RED"
)

type batteryLevelIcon string

const (
	BATTERY_LEVEL_4_4 batteryLevelIcon = "\uf240"
	BATTERY_LEVEL_3_4 batteryLevelIcon = "\uf241"
	BATTERY_LEVEL_2_4 batteryLevelIcon = "\uf242"
	BATTERY_LEVEL_1_4 batteryLevelIcon = "\uf243"
	BATTERY_LEVEL_0_4 batteryLevelIcon = "\uf244"
)

type ResponseRenderer struct {
	template     core.Template
	nodeId       string
	itemRenderer []core.Renderer
}

type responseData struct {
	RenderId string
	NodeId   string
	Items    string
}

type ErrorRenderer struct {
	template core.Template
	err      error
}

type billingReportData struct {
	Anchor core.Point
	Period string
	Amount string
}

type billingReportAmount struct {
	Amount   float64
	Currency string
}

type BillingReportRenderer struct {
	template        core.Template
	anchor          core.Point
	logger          log.Logger
	reportCurrency  string
	displayCurrency string
	datasource      core.DataSource
	dataSourceChan  <-chan proto.Message
	billingReport   *billingReportData
	exchangeRates   map[string]*events.ExchangeRate
}

type WeatherRenderer struct {
	currentWeatherTemplate core.Template
	forecastTemplate       core.Template
	anchor                 core.Point
	logger                 log.Logger
	datasource             core.DataSource
	dataSourceChan         <-chan proto.Message
	weatherData            *events.WeatherData
	weatherIconMap         WeatherIconMap
}

type weatherData struct {
	Anchor       core.Point
	WeatherIcon  string
	Temperature  string
	WindSpeed    string
	Day          string
	DisplayIndex int
}

type WeatherIconMap struct {
	icons map[string]string
}

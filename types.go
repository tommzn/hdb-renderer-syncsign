package syncsign

import (
	"github.com/golang/protobuf/proto"
	log "github.com/tommzn/go-log"
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

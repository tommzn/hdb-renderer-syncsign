package syncsign

import (
	"context"
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	core "github.com/tommzn/hdb-renderer-core"
)

// NewIndoorClimateRenderer returns a new renderer for infoor climate data. Room will be taken from passed config, template and datasource have to be passed.
func NewIndoorClimateRenderer(conf config.Config, logger log.Logger, template core.Template, datasource core.DataSource) *IndoorClimateRenderer {

	anchor := anchorFromConfig(conf, "hdb.indoorclimate.anchor")
	fmt.Printf("Anchor, X: %d, Y: %d\n", anchor.X, anchor.Y)
	size := sizeFromConfig(conf, "hdb.indoorclimate.size")
	spacing := spacingFromConfig(conf, "hdb.indoorclimate.spacing")
	roomCfg := configForRooms(conf, "hdb.indoorclimate")
	return &IndoorClimateRenderer{
		originAnchor: anchor,
		size:         size,
		spacing:      spacing,
		datasource:   datasource,
		template:     template,
		logger:       logger,
		roomClimate:  make(map[string]indoorCliemate),
		roomCfg:      roomCfg,
	}
}

// Content fetches current inddor climate data and generated room climate elements based
// pn given room/device config.
func (renderer *IndoorClimateRenderer) Content() (string, error) {

	defer renderer.logger.Flush()

	if len(renderer.roomCfg.rooms) == 0 {
		renderer.logger.Error("No room config for rendering!")
		return "", nil
	}

	// Init indoor climate data from used datasource if nothing is available
	// or if renderer doesn't observer datasource actively.
	if len(renderer.roomClimate) == 0 || renderer.dataSourceChan == nil {
		renderer.initIndoorClimateData()
	}

	if len(renderer.roomClimate) == 0 {
		renderer.logger.Error("No room climate to render.")
		return "", nil
	}

	content := ""
	anchor := renderer.originAnchor
	roomClimate := renderer.sortedRoomClimateData()
	for _, climate := range roomClimate {
		climate.Anchor = anchor
		elementContent, err := renderer.template.RenderWith(climate)
		if err != nil {
			return "", err
		}
		content = content + elementContent
		anchor.X = anchor.X + renderer.size.Width + renderer.spacing.Left + renderer.spacing.Right
	}
	return content, nil
}

// ObserveDataSource will listen for new indoor climate data provided by used datasource.
func (renderer *IndoorClimateRenderer) ObserveDataSource(ctx context.Context) {

	defer renderer.logger.Flush()

	filter := []hdbcore.DataSource{hdbcore.DATASOURCE_INDOORCLIMATE}
	renderer.dataSourceChan = renderer.datasource.Observe(&filter)
	for {
		select {
		case message, ok := <-renderer.dataSourceChan:
			if !ok {
				renderer.logger.Error("Error at reading datasource channel. Stop observing!")
				return
			}
			renderer.addAsIndoorClimateData(message)
		case <-ctx.Done():
			renderer.logger.Info("Camceled, stop observing.")
			return
		}
	}
}

// InitIndoorClimateData will dop existing indoor climate data and fetch all available events from used datasource.
func (renderer *IndoorClimateRenderer) initIndoorClimateData() {

	renderer.roomClimate = make(map[string]indoorCliemate)

	messages, err := renderer.datasource.All(hdbcore.DATASOURCE_INDOORCLIMATE)
	if err != nil {
		renderer.logger.Error("Unable to get indoor climate, reason: ", err)
		return
	}
	renderer.logger.Infof("Fetch %d indoor climate messages", len(messages))

	for _, message := range messages {
		renderer.addAsIndoorClimateData(message)
	}
}

// addAsIndoorClimateData will try to add passed message to local indoor climate data.
func (renderer *IndoorClimateRenderer) addAsIndoorClimateData(message proto.Message) {

	if indoorClimate, ok := message.(*events.IndoorClimate); ok {
		if roomId, ok := renderer.roomCfg.deviceMap[indoorClimate.DeviceId]; ok {
			roomClimate := renderer.getRoomClimate(roomId)
			switch indoorClimate.Type {
			case events.MeasurementType_TEMPERATURE:
				roomClimate.Temperature = formatTemperature(indoorClimate.Value)
			case events.MeasurementType_HUMIDITY:
				roomClimate.Humidity = formatHumidity(indoorClimate.Value)
			case events.MeasurementType_BATTERY:
				roomClimate.BatteryIcon = batteryIcon(indoorClimate.Value)
				roomClimate.BatteryIconColor = batteryIconColor(indoorClimate.Value)
			}
			renderer.roomClimate[roomId] = roomClimate
		}
	}
}

// getRoomClimate will have a look if there's already climate data for given room.
// If nothing exists a new room climate with default values will created, assigned to passed
// room and returned.
func (renderer *IndoorClimateRenderer) getRoomClimate(roomId string) indoorCliemate {

	if roomClimate, ok := renderer.roomClimate[roomId]; ok {
		return roomClimate
	}

	roomClimate := indoorCliemate{
		DisplayIndex:     "0",
		Temperature:      "--",
		Humidity:         "--",
		BatteryIcon:      BATTERY_LEVEL_0_4,
		BatteryIconColor: COLOR_BLACK,
		RoomName:         "Room",
		Anchor:           core.Point{X: 0, Y: 0},
	}
	if roomCfg, ok := renderer.roomCfg.rooms[roomId]; ok {
		roomClimate.DisplayIndex = roomCfg.DisplayIndex
		roomClimate.RoomName = roomCfg.Name
	}
	return roomClimate
}

// sortedRoomClimateData sorts current room climate based on displayIndex given by room config.
func (renderer *IndoorClimateRenderer) sortedRoomClimateData() []indoorCliemate {

	roomClimate := []indoorCliemate{}
	for _, cliamte := range renderer.roomClimate {
		roomClimate = append(roomClimate, cliamte)
	}
	sort.Slice(roomClimate, func(i, j int) bool {
		return roomClimate[i].DisplayIndex < roomClimate[j].DisplayIndex
	})
	return roomClimate
}

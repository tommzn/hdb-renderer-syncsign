package syncsign

import (
	"fmt"
	"strconv"
	"strings"

	config "github.com/tommzn/go-config"
	core "github.com/tommzn/hdb-renderer-core"
)

// anchorFromConfig will try to get an anchor defined in passed config.
// If there's no anchor defined, or if it's incomplete an anchor at X:0, Y:0 will be returned.
// Negative values will be used as 0.
func anchorFromConfig(conf config.Config, anchorConfigKey string) core.Point {

	positionY := conf.GetAsInt(anchorConfigKey+".y", config.AsIntPtr(0))
	positionX := conf.GetAsInt(anchorConfigKey+".x", config.AsIntPtr(0))
	return core.Point{X: forcePositive(*positionX), Y: forcePositive(*positionY)}
}

// sizeConfigKey will try to load size for given config key and returns default
// values of Height: 0 and Width: 0 if no cofig exists.
// Negative values will be used as 0.
func sizeFromConfig(conf config.Config, sizeConfigKey string) core.Size {

	height := conf.GetAsInt(sizeConfigKey+".height", config.AsIntPtr(0))
	width := conf.GetAsInt(sizeConfigKey+".width", config.AsIntPtr(0))
	return core.Size{Height: forcePositive(*height), Width: forcePositive(*width)}
}

// spacingFromConfig will try to load distance config between elements and returns
// spacing of 0 if nothing has been found in passed config.
// Negative values will be used as 0.
func spacingFromConfig(conf config.Config, spacingConfigKey string) core.Spacing {

	all := conf.GetAsInt(spacingConfigKey, nil)
	top := conf.GetAsInt(spacingConfigKey+".top", config.AsIntPtr(0))
	left := conf.GetAsInt(spacingConfigKey+".left", config.AsIntPtr(0))
	right := conf.GetAsInt(spacingConfigKey+".right", config.AsIntPtr(0))
	bottom := conf.GetAsInt(spacingConfigKey+".bottom", config.AsIntPtr(0))

	if all != nil && *all != 0 {
		*all = forcePositive(*all)
		return core.Spacing{Top: *all, Left: *all, Right: *all, Bottom: *all}
	}
	return core.Spacing{Top: forcePositive(*top), Left: forcePositive(*left), Right: forcePositive(*right), Bottom: forcePositive(*bottom)}
}

// configForRooms will try to extract room and device settings from passed config.
func configForRooms(conf config.Config, configKey string) roomConfig {

	roomsCfg := roomConfig{
		rooms:     make(map[string]room),
		deviceMap: make(map[string]string),
	}

	rooms := conf.GetAsSliceOfMaps(configKey + ".rooms")
	for idx, roomCfg := range rooms {
		if roomId, ok := roomCfg["id"]; ok {
			roomName, ok1 := roomCfg["name"]
			if !ok1 {
				roomName = fmt.Sprintf("Room %d", idx)
			}
			displayIndex, ok1 := roomCfg["displayIndex"]
			if !ok1 {
				displayIndex = fmt.Sprintf("%d", idx)
			}
			roomsCfg.rooms[roomId] = room{
				Id:           roomId,
				Name:         roomName,
				DisplayIndex: displayIndex,
			}
		}
	}

	devices := conf.GetAsSliceOfMaps(configKey + ".devices")
	for _, device := range devices {
		if id, ok := device["id"]; ok {
			if roomId, ok1 := device["roomId"]; ok1 {
				roomsCfg.deviceMap[id] = roomId
			}
		}
	}
	return roomsCfg
}

// forcePositive will return 0 for all negative values and origin value for all others.
func forcePositive(val int) int {
	if val < 0 {
		return 0
	}
	return val
}

func formatTemperature(temperature string) string {
	if floatTemp, err := strconv.ParseFloat(temperature, 64); err == nil {
		return fmt.Sprintf("%.1f", floatTemp)
	}
	return temperature
}

func formatHumidity(humidity string) string {
	if floatHum, err := strconv.ParseFloat(humidity, 64); err == nil {
		return fmt.Sprintf("%.0f", floatHum)
	}
	return humidity
}

func batteryIcon(batteryValue string) batteryLevelIcon {
	if intVal, err := strconv.Atoi(batteryValue); err == nil {
		switch {
		case intVal >= 90:
			return BATTERY_LEVEL_4_4
		case intVal >= 75:
			return BATTERY_LEVEL_3_4
		case intVal >= 50:
			return BATTERY_LEVEL_2_4
		case intVal >= 10:
			return BATTERY_LEVEL_1_4
		case intVal <= 10:
			return BATTERY_LEVEL_0_4
		}
	}
	return BATTERY_LEVEL_0_4
}

func batteryIconColor(batteryValue string) textColor {
	if intVal, err := strconv.Atoi(batteryValue); err == nil && intVal <= 5 {
		return COLOR_RED
	}
	return COLOR_BLACK
}

// AppendItems appends passed items to given content, separated by default JSON element separattor: ",".
// Leasing separators in content, or trailing separators in items will be removed.
func appendItems(items, newItems string) string {
	fmt.Println(items)
	fmt.Println(newItems)
	if newItems != "" {
		newItems = strings.TrimPrefix(newItems, ",")
		fmt.Println(newItems)
		newItems = strings.TrimSuffix(newItems, ",")
		fmt.Println(newItems)
		items = strings.TrimPrefix(items, ",")
		fmt.Println(items)
		items = strings.TrimSuffix(items, ",")
		fmt.Println(items)
		items = items + "," + newItems
		fmt.Println(items)
	}
	fmt.Println(items)
	return strings.TrimPrefix(strings.TrimSuffix(items, ","), ",")
}

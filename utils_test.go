package syncsign

import (
	"github.com/stretchr/testify/suite"
	"testing"

	config "github.com/tommzn/go-config"
	core "github.com/tommzn/hdb-renderer-core"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestGetAncharFromConfig() {

	testCases := map[string]core.Point{
		"fixtures/anchortest01.yml": core.Point{X: 100, Y: 200},
		"fixtures/anchortest02.yml": core.Point{X: 0, Y: 200},
		"fixtures/anchortest03.yml": core.Point{X: 100, Y: 0},
		"fixtures/anchortest04.yml": core.Point{X: 0, Y: 0},
		"fixtures/anchortest05.yml": core.Point{X: 0, Y: 0},
	}
	for configFile, expectedAnchor := range testCases {
		conf := loadConfigForTest(config.AsStringPtr(configFile))
		anchor := anchorFromConfig(conf, "hdb.indoorclimate.anchor")
		suite.Equal(expectedAnchor, anchor)
	}
}

func (suite *UtilsTestSuite) TestGetSizeFromConfig() {

	testCases := map[string]core.Size{
		"fixtures/sizetest01.yml": core.Size{Height: 100, Width: 200},
		"fixtures/sizetest02.yml": core.Size{Height: 0, Width: 200},
		"fixtures/sizetest03.yml": core.Size{Height: 100, Width: 0},
		"fixtures/sizetest04.yml": core.Size{Height: 0, Width: 0},
		"fixtures/sizetest05.yml": core.Size{Height: 0, Width: 0},
	}
	for configFile, expectedSize := range testCases {
		conf := loadConfigForTest(config.AsStringPtr(configFile))
		size := sizeFromConfig(conf, "hdb.indoorclimate.size")
		suite.Equal(expectedSize, size)
	}
}

func (suite *UtilsTestSuite) TestGetSpacingFromConfig() {

	testCases := map[string]core.Spacing{
		"fixtures/spacingtest01.yml": core.Spacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
		"fixtures/spacingtest02.yml": core.Spacing{Top: 10, Left: 0, Right: 0, Bottom: 0},
		"fixtures/spacingtest03.yml": core.Spacing{Top: 0, Left: 10, Right: 0, Bottom: 0},
		"fixtures/spacingtest04.yml": core.Spacing{Top: 0, Left: 0, Right: 10, Bottom: 0},
		"fixtures/spacingtest05.yml": core.Spacing{Top: 0, Left: 0, Right: 0, Bottom: 10},
		"fixtures/spacingtest06.yml": core.Spacing{Top: 0, Left: 0, Right: 0, Bottom: 0},
		"fixtures/spacingtest07.yml": core.Spacing{Top: 0, Left: 0, Right: 0, Bottom: 0},
		"fixtures/spacingtest08.yml": core.Spacing{Top: 0, Left: 0, Right: 0, Bottom: 0},
	}
	for configFile, expectedSpacing := range testCases {
		conf := loadConfigForTest(config.AsStringPtr(configFile))
		spacing := spacingFromConfig(conf, "hdb.indoorclimate.spacing")
		suite.Equal(expectedSpacing, spacing)
	}
}

func (suite *UtilsTestSuite) TestGetRoomConfig() {

	conf := loadConfigForTest(nil)
	roomCfg := configForRooms(conf, "hdb.indoorclimate")
	suite.Len(roomCfg.rooms, 4)
	suite.Len(roomCfg.deviceMap, 3)
}

func (suite *UtilsTestSuite) TestFormatValues() {

	suite.Equal("23.4", formatTemperature("23.4"))
	suite.Equal("23.4", formatTemperature("23.42"))
	suite.Equal("23.5", formatTemperature("23.47"))
	suite.Equal("23.0", formatTemperature("23"))
	suite.Equal("xxx", formatTemperature("xxx"))

	suite.Equal("23", formatHumidity("23.4"))
	suite.Equal("23", formatHumidity("23.42"))
	suite.Equal("24", formatHumidity("23.67"))
	suite.Equal("23", formatHumidity("23"))
	suite.Equal("xxx", formatHumidity("xxx"))
}

func (suite *UtilsTestSuite) TestParseBatteryValue() {

	intVal01 := batteryValueToInt("96.7")
	suite.Equal(96, intVal01)
	intVal02 := batteryValueToInt("93")
	suite.Equal(93, intVal02)
	intVal03 := batteryValueToInt("xxx")
	suite.Equal(0, intVal03)
}

func (suite *UtilsTestSuite) TestBatteryIcon() {

	suite.Equal(COLOR_BLACK, batteryIconColor("100"))
	suite.Equal(COLOR_BLACK, batteryIconColor("10"))
	suite.Equal(COLOR_RED, batteryIconColor("5"))
	suite.Equal(COLOR_RED, batteryIconColor("1"))
	suite.Equal(COLOR_RED, batteryIconColor("4"))
	suite.Equal(COLOR_RED, batteryIconColor("0"))
	suite.Equal(COLOR_RED, batteryIconColor("xxx"))

	suite.Equal(BATTERY_LEVEL_0_4, batteryIcon("xxx"))
	suite.Equal(BATTERY_LEVEL_4_4, batteryIcon("100"))
	suite.Equal(BATTERY_LEVEL_4_4, batteryIcon("90"))
	suite.Equal(BATTERY_LEVEL_3_4, batteryIcon("89"))
	suite.Equal(BATTERY_LEVEL_3_4, batteryIcon("80"))
	suite.Equal(BATTERY_LEVEL_3_4, batteryIcon("75"))
	suite.Equal(BATTERY_LEVEL_2_4, batteryIcon("70"))
	suite.Equal(BATTERY_LEVEL_2_4, batteryIcon("60"))
	suite.Equal(BATTERY_LEVEL_2_4, batteryIcon("50"))
	suite.Equal(BATTERY_LEVEL_1_4, batteryIcon("40"))
	suite.Equal(BATTERY_LEVEL_1_4, batteryIcon("30"))
	suite.Equal(BATTERY_LEVEL_1_4, batteryIcon("20"))
	suite.Equal(BATTERY_LEVEL_1_4, batteryIcon("10"))
	suite.Equal(BATTERY_LEVEL_0_4, batteryIcon("5"))
}

func (suite *UtilsTestSuite) TestConvertDegreesToDirection() {

	suite.Equal("N", degreesToDirection(0))
	suite.Equal("N", degreesToDirection(360))
	suite.Equal("N", degreesToDirection(20))
	suite.Equal("NE", degreesToDirection(30))
	suite.Equal("E", degreesToDirection(87))
	suite.Equal("SE", degreesToDirection(120))
	suite.Equal("S", degreesToDirection(175))
	suite.Equal("SW", degreesToDirection(220))
	suite.Equal("W", degreesToDirection(271))
	suite.Equal("NW", degreesToDirection(320))
	suite.Equal("N", degreesToDirection(340))
	suite.Equal("N/A", degreesToDirection(600))
}

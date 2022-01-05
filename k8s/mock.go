package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	core "github.com/tommzn/hdb-renderer-core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type dataSourceMock struct {
	indoorClimateDevices []string
	measurementTypes     []events.MeasurementType
	currentValues        map[events.MeasurementType]float64
	publishInterval      time.Duration
	stackSize            int
	events               map[hdbcore.DataSource][]proto.Message
	eventChan            chan proto.Message
	chanFilter           []hdbcore.DataSource
}

func newDataSourceMock(indoorClimateDevices []string) core.DataSource {

	return &dataSourceMock{
		indoorClimateDevices: indoorClimateDevices,
		measurementTypes: []events.MeasurementType{
			events.MeasurementType_TEMPERATURE,
			events.MeasurementType_HUMIDITY,
			events.MeasurementType_BATTERY,
		},
		currentValues:   make(map[events.MeasurementType]float64),
		publishInterval: 10 * time.Second,
		stackSize:       100,
		events:          make(map[hdbcore.DataSource][]proto.Message),
		eventChan:       make(chan proto.Message, 100),
		chanFilter:      []hdbcore.DataSource{},
	}
}

func (mock *dataSourceMock) Run(ctx context.Context) {

	for {
		select {
		case <-time.After(mock.publishInterval):
			mock.publishNewMessage()
		case <-ctx.Done():
			return
		}
	}
}

func (mock *dataSourceMock) initMessages() {

	for _, deviceId := range mock.indoorClimateDevices {
		for _, measurementType := range mock.measurementTypes {

			currentValue, ok := mock.currentValues[measurementType]
			if !ok {
				currentValue = defauktValueForMeasurementType(measurementType)
			}
			message := &events.IndoorClimate{
				Timestamp: timestamppb.New(time.Now()),
				DeviceId:  deviceId,
				Type:      measurementType,
				Value:     formatValue(measurementType, currentValue),
			}
			mock.appendToStack(message, hdbcore.DATASOURCE_INDOORCLIMATE)
			mock.writeToChannel(message, hdbcore.DATASOURCE_INDOORCLIMATE)
		}
	}
}

func (mock *dataSourceMock) publishNewMessage() {

	measurementType := mock.randomSelectMeasurementType()
	deviceId := mock.randomSelectDeviceId()

	currentValue, ok := mock.currentValues[measurementType]
	if !ok {
		currentValue = defauktValueForMeasurementType(measurementType)
	}
	newValue := newValueForMeasurementType(measurementType, currentValue)
	mock.currentValues[measurementType] = newValue
	message := &events.IndoorClimate{
		Timestamp: timestamppb.New(time.Now()),
		DeviceId:  deviceId,
		Type:      measurementType,
		Value:     formatValue(measurementType, newValue),
	}
	mock.appendToStack(message, hdbcore.DATASOURCE_INDOORCLIMATE)
	mock.writeToChannel(message, hdbcore.DATASOURCE_INDOORCLIMATE)
}

func (mock *dataSourceMock) randomSelectMeasurementType() events.MeasurementType {
	rand.Seed(time.Now().UnixNano())
	return mock.measurementTypes[rand.Intn(len(mock.measurementTypes))]
}

func (mock *dataSourceMock) randomSelectDeviceId() string {
	rand.Seed(time.Now().UnixNano())
	return mock.indoorClimateDevices[rand.Intn(len(mock.indoorClimateDevices))]
}

func (mock *dataSourceMock) appendToStack(message proto.Message, datasource hdbcore.DataSource) {
	if events, ok := mock.events[hdbcore.DATASOURCE_INDOORCLIMATE]; ok {
		if len(events) == mock.stackSize {
			events = events[1:]
		}
		mock.events[hdbcore.DATASOURCE_INDOORCLIMATE] = append(events, message)
	} else {
		mock.events[hdbcore.DATASOURCE_INDOORCLIMATE] = []proto.Message{message}
	}
}

func (mock *dataSourceMock) writeToChannel(message proto.Message, datasource hdbcore.DataSource) {
	if mock.isInFilter(datasource) &&
		len(mock.eventChan) < cap(mock.eventChan) {
		mock.eventChan <- message
	}
}

func (mock *dataSourceMock) isInFilter(datasource hdbcore.DataSource) bool {

	if len(mock.chanFilter) == 0 {
		return true
	}

	for _, filterItem := range mock.chanFilter {
		if datasource == filterItem {
			return true
		}
	}
	return false
}

func (mock *dataSourceMock) Latest(datasource hdbcore.DataSource) (proto.Message, error) {
	if events, ok := mock.events[datasource]; ok && len(events) > 0 {
		return events[len(events)-1], nil
	}
	return nil, errors.New("No events available.")
}

func (mock *dataSourceMock) All(datasource hdbcore.DataSource) ([]proto.Message, error) {
	if events, ok := mock.events[datasource]; ok && len(events) > 0 {
		return events, nil
	}
	return nil, errors.New("No events available.")
}

func (mock *dataSourceMock) Observe(filter *[]hdbcore.DataSource) <-chan proto.Message {
	if filter != nil {
		mock.chanFilter = *filter
	}
	return mock.eventChan
}

func formatValue(measurementType events.MeasurementType, value float64) string {
	switch measurementType {
	case events.MeasurementType_TEMPERATURE:
		return fmt.Sprintf("%.1f", value)
	case events.MeasurementType_HUMIDITY:
		return fmt.Sprintf("%.0f", value)
	default: // e.g. events.MeasurementType_BATTERY
		return fmt.Sprintf("%d", int(value))
	}
}

func newValueForMeasurementType(measurementType events.MeasurementType, currentValue float64) float64 {

	rand.Seed(time.Now().UnixNano())
	step := 0.0
	switch measurementType {
	case events.MeasurementType_HUMIDITY:
		if rand.Intn(100) > 50 {
			step = 1.0
		} else {
			step = -1.0
		}
	case events.MeasurementType_BATTERY:
		step = -1.0
	default: // e.g. events.MeasurementType_TEMPERATURE
		if rand.Intn(100) > 50 {
			step = 0.1
		} else {
			step = -0.1
		}
	}
	return currentValue + step
}

func defauktValueForMeasurementType(measurementType events.MeasurementType) float64 {
	switch measurementType {
	case events.MeasurementType_TEMPERATURE:
		return 23.5
	case events.MeasurementType_HUMIDITY:
		return 64.0
	default: // e.g. events.MeasurementType_BATTERY
		return 100.0
	}
}
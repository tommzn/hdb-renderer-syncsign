package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/tommzn/go-log"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	core "github.com/tommzn/hdb-renderer-core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type dataSourceMock struct {
	logger               log.Logger
	indoorClimateDevices []string
	measurementTypes     []events.MeasurementType
	currentValues        map[string]map[events.MeasurementType]float64
	publishInterval      time.Duration
	stackSize            int
	events               map[hdbcore.DataSource][]proto.Message
	eventChan            chan proto.Message
	chanFilter           []hdbcore.DataSource
}

func newDataSourceMock(indoorClimateDevices []string, logger log.Logger) core.DataSource {

	return &dataSourceMock{
		logger:               logger,
		indoorClimateDevices: indoorClimateDevices,
		measurementTypes: []events.MeasurementType{
			events.MeasurementType_TEMPERATURE,
			events.MeasurementType_HUMIDITY,
			events.MeasurementType_BATTERY,
		},
		currentValues:   make(map[string]map[events.MeasurementType]float64),
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

			currentValue, ok := mock.currentValues[deviceId][measurementType]
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
	mock.publisBillingReport()
	mock.publisExchangeRate()
}

func (mock *dataSourceMock) publishNewMessage() {

	measurementType := mock.randomSelectMeasurementType()
	deviceId := mock.randomSelectDeviceId()

	currentValue, ok := mock.currentValues[deviceId][measurementType]
	if !ok {
		mock.currentValues[deviceId] = make(map[events.MeasurementType]float64)
		currentValue = defauktValueForMeasurementType(measurementType)
	}
	newValue := newValueForMeasurementType(measurementType, currentValue)
	mock.currentValues[deviceId][measurementType] = newValue
	message := &events.IndoorClimate{
		Timestamp: timestamppb.New(time.Now()),
		DeviceId:  deviceId,
		Type:      measurementType,
		Value:     formatValue(measurementType, newValue),
	}
	mock.appendToStack(message, hdbcore.DATASOURCE_INDOORCLIMATE)
	mock.writeToChannel(message, hdbcore.DATASOURCE_INDOORCLIMATE)

	mock.publisBillingReport()
	mock.publisExchangeRate()
}

func (mock *dataSourceMock) publisBillingReport() {

	billingAmount := make(map[string]float64)
	taxAmount := make(map[string]float64)
	billingAmount["xxx"] = 5.14
	billingAmount["zzz"] = 12.53
	taxAmount["xxx"] = 0.87
	taxAmount["zzz"] = 2.15
	billingReport := &events.BillingReport{
		BillingPeriod: "Jan 2022",
		BillingAmount: billingAmount,
		TaxAmount:     taxAmount,
	}
	mock.appendToStack(billingReport, hdbcore.DATASOURCE_BILLINGREPORT)
	mock.writeToChannel(billingReport, hdbcore.DATASOURCE_BILLINGREPORT)
}

func (mock *dataSourceMock) publisExchangeRate() {

	exchangeRate := &events.ExchangeRate{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Rate:         0.8345,
		Timestamp:    timestamppb.New(time.Now()),
	}
	mock.appendToStack(exchangeRate, hdbcore.DATASOURCE_EXCHANGERATE)
	mock.writeToChannel(exchangeRate, hdbcore.DATASOURCE_EXCHANGERATE)
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
	if events, ok := mock.events[datasource]; ok {
		if len(events) == mock.stackSize {
			events = events[1:]
		}
		mock.events[datasource] = append(events, message)
	} else {
		mock.events[datasource] = []proto.Message{message}
	}
}

func (mock *dataSourceMock) writeToChannel(message proto.Message, datasource hdbcore.DataSource) {
	mock.logger.Debugf("Publish new event to %s", datasource)
	if mock.isInFilter(datasource) &&
		len(mock.eventChan) < cap(mock.eventChan) {
		mock.eventChan <- message
		return
	}
	mock.logger.Debugf("No subscription for datasource or chanel blocked %d/%d", datasource, len(mock.eventChan), cap(mock.eventChan))
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
	mock.appendToDataSourceFilter(filter)
	return mock.eventChan
}

func (mock *dataSourceMock) appendToDataSourceFilter(filter *[]hdbcore.DataSource) {
	if filter != nil {
		if mock.chanFilter == nil || len(mock.chanFilter) == 0 {
			mock.chanFilter = *filter
		} else {
			for _, datasource := range *filter {
				mock.chanFilter = append(mock.chanFilter, datasource)
			}
		}
	}
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
	steps := []float64{}
	switch measurementType {
	case events.MeasurementType_HUMIDITY:
		steps = []float64{1.0, -1.0, 2.0, -2.0, 3.0, -3.0}
	case events.MeasurementType_BATTERY:
		steps = []float64{0.0, -1.0, -2.0}
	default: // e.g. events.MeasurementType_TEMPERATURE
		steps = []float64{0.1, -0.1, 0.2, -0.2, 0.3, -0.3}
	}
	rand.Seed(time.Now().UnixNano())
	return currentValue + steps[rand.Int()%len(steps)]
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

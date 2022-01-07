package syncsign

import (
	"context"
	"errors"

	"github.com/golang/protobuf/proto"
	hdbcore "github.com/tommzn/hdb-core"
	core "github.com/tommzn/hdb-renderer-core"
)

type datasourceMock struct {
	shouldReturnError bool
	shouldReturnEmpty bool
	data              map[hdbcore.DataSource][]proto.Message
	eventChan         chan proto.Message
}

func newDataSourceMock(shouldReturnError, shouldReturnEmpty bool, data map[hdbcore.DataSource][]proto.Message) core.DataSource {
	chanLen := 0
	for _, messages := range data {
		chanLen += len(messages)
	}
	return &datasourceMock{
		shouldReturnError: shouldReturnError,
		shouldReturnEmpty: shouldReturnError,
		data:              data,
		eventChan:         make(chan proto.Message, chanLen+10),
	}
}

func (mock *datasourceMock) Latest(datasource hdbcore.DataSource) (proto.Message, error) {

	if mock.shouldReturnError || mock.shouldReturnEmpty {
		return nil, errors.New("Error occured!")
	}
	return mock.data[datasource][len(mock.data[datasource])-1], nil
}

func (mock *datasourceMock) All(datasource hdbcore.DataSource) ([]proto.Message, error) {

	events := []proto.Message{}
	if mock.shouldReturnError {
		return events, errors.New("Error occured!")
	}
	if mock.shouldReturnEmpty {
		return events, nil
	}
	return mock.data[datasource], nil
}

func (mock *datasourceMock) Observe(filter *[]hdbcore.DataSource) <-chan proto.Message {

	if filter != nil {
		for _, datasource := range *filter {
			for _, message := range mock.data[datasource] {
				mock.eventChan <- message
			}
		}
	}
	for _, messages := range mock.data {
		for _, message := range messages {
			mock.eventChan <- message
		}
	}
	return mock.eventChan
}

func (mock *datasourceMock) writeToMessageChannel(message proto.Message) {
	mock.eventChan <- message
}

type failingTemplate struct {
}

func newFailingTemplate() core.Template {
	return &failingTemplate{}
}

func (mock *failingTemplate) RenderWith(interface{}) (string, error) {
	return "", errors.New("Error occured!")
}

type rendererMock struct {
	shouldReturnEmptyContent bool
	shouldFail               bool
}

func (renderer *rendererMock) Size() core.Size {
	return core.Size{Height: 0, Width: 0}
}

func newRendererMock(shouldReturnEmptyContent, shouldFail bool) core.Renderer {
	return &rendererMock{
		shouldReturnEmptyContent: shouldReturnEmptyContent,
		shouldFail:               shouldFail,
	}
}

func (renderer *rendererMock) Content() (string, error) {

	if renderer.shouldFail {
		return "", errors.New("Error occured!")
	}

	if renderer.shouldReturnEmptyContent {
		return "", nil
	} else {
		return "{\"id\": \"Item1\"},{\"id\": \"Item2\"},", nil
	}
}

func (renderer *rendererMock) ObserveDataSource(ctx context.Context) {

}

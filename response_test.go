package syncsign

import (
	"context"
	"github.com/stretchr/testify/suite"
	core "github.com/tommzn/hdb-renderer-core"
	"testing"
)

type ResponseTestSuite struct {
	suite.Suite
}

func TestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}

func (suite *ResponseTestSuite) TestGenerateContent() {

	itemRenderer := []core.Renderer{
		newRendererMock(true, false),
		newRendererMock(false, false),
		newRendererMock(true, false),
		newRendererMock(false, false),
	}
	tmpl := templateQithFileForTest("templates/response.json")
	renderer := NewResponseRenderer(tmpl, "Node-1", itemRenderer)

	content, err := renderer.Content()
	suite.Nil(err)

	// Replace renderer id and timestamp with default value for assertion
	content = replaceUUID(content, "RenderId-1")
	content = replaceTimeStamp(content, "TimeStamp-1")
	assertTemplateHash(suite.Assert(), content, "03729f689b446819bb1ef11ca42517e87bcf5cf3")

	size := renderer.Size()
	suite.Equal(528, size.Height)
	suite.Equal(880, size.Width)

	renderer.ObserveDataSource(context.Background())
}

func (suite *ResponseTestSuite) TestGenerateContentWithoutItemRenderer() {

	tmpl := templateQithFileForTest("templates/response.json")
	renderer := NewResponseRenderer(tmpl, "Node-1", []core.Renderer{})

	content, err := renderer.Content()
	suite.NotNil(err)
	suite.Equal("", content)
}

func (suite *ResponseTestSuite) TestGenerateContentWithFailingItemRenderer() {

	itemRenderer := []core.Renderer{
		newRendererMock(true, false),
		newRendererMock(false, true),
		newRendererMock(false, false),
	}
	tmpl := templateQithFileForTest("templates/response.json")
	renderer := NewResponseRenderer(tmpl, "Node-1", itemRenderer)

	content, err := renderer.Content()
	suite.NotNil(err)
	suite.Equal("", content)
}

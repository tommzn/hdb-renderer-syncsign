package syncsign

import (
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
	assertTemplateHash(suite.Assert(), content, "b38fe380db159d013db97d4b601609521aad9bfd")
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

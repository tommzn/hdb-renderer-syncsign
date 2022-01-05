package syncsign

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorRendererTestSuite struct {
	suite.Suite
}

func TestErrorRendererTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorRendererTestSuite))
}

func (suite *ErrorRendererTestSuite) TestGenerateContent() {

	tmpl := templateQithFileForTest("templates/error.json")
	err := errors.New("Failed to generate content, Err 101")
	renderer := NewErrorRenderer(tmpl, err)

	content, err := renderer.Content()
	suite.Nil(err)

	// Replace renderer id and timestamp with default value for assertion
	content = replaceUUID(content, "RenderId-1")
	content = replaceTimeStamp(content, "TimeStamp-1")
	assertTemplateHash(suite.Assert(), content, "826eb61847fb54742407c6a979110ec84f276213")
}

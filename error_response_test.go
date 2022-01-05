package syncsign

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorResponseTestSuite struct {
	suite.Suite
}

func TestErrorResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorResponseTestSuite))
}

func (suite *ErrorResponseTestSuite) TestGenerateContent() {

	tmpl := templateQithFileForTest("templates/error_response.json")
	err := errors.New("Failed to generate content, Err 101")
	renderer := NewErrorRenderer(tmpl, "Node-1", err)

	content, err := renderer.Content()
	suite.Nil(err)

	// Replace renderer id and timestamp with default value for assertion
	content = replaceUUID(content, "RenderId-1")
	content = replaceTimeStamp(content, "TimeStamp-1")
	assertTemplateHash(suite.Assert(), content, "25e3928425a58ebeec9317c6f239f3592f6918cb")

	size := renderer.Size()
	suite.Equal(124, size.Height)
	suite.Equal(400, size.Width)

	renderer.ObserveDataSource(context.Background())
}

package syncsign

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TimestampTestSuite struct {
	suite.Suite
}

func TestTimestampTestSuite(t *testing.T) {
	suite.Run(t, new(TimestampTestSuite))
}

func (suite *TimestampTestSuite) TestGenerateContent() {

	tmpl := templateQithFileForTest("templates/timestamp.json")
	renderer := NewTimestampRenderer(tmpl)

	content, err := renderer.Content()
	suite.Nil(err)
	content = replaceTimeStamp(content, "TimeStamp-1")
	assertTemplateHash(suite.Assert(), content, "62e2a622de3f26945696c8649154cbf272a01327")
}

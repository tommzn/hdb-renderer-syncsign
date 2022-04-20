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
	assertTemplateHash(suite.Assert(), content, "5b204be4b22d306b9bfc79eb05afaa306f68c912")
}

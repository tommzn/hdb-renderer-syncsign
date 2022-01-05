package syncsign

import (
	"time"

	core "github.com/tommzn/hdb-renderer-core"
)

// NewTimestampRenderer returns a new renderer to generate a single item which contains a timestamp.
func NewTimestampRenderer(template core.Template) core.Renderer {
	return &TimestampRenderer{
		template: template,
	}
}

// Content generates a single item with a current timestamp.
func (renderer *TimestampRenderer) Content() (string, error) {
	return renderer.template.RenderWith(time.Now().Format("2006-01-02 15:04:05 MST"))
}

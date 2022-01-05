package syncsign

import (
	core "github.com/tommzn/hdb-renderer-core"
)

// NewErrorRenderertemplate returns a renderer which generates items for passed error message.
func NewErrorRenderer(template core.Template, err error) core.Renderer {
	return &ErrorRenderer{
		template: template,
		err:      err,
	}
}

// Content returns errpr message passed at initialization as items.
// Together with a title and a timestamp.
func (renderer *ErrorRenderer) Content() (string, error) {
	return renderer.template.RenderWith(renderer.err.Error())
}

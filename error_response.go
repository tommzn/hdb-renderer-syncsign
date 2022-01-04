package syncsign

import (
	"context"
	"time"

	utils "github.com/tommzn/go-utils"
	core "github.com/tommzn/hdb-renderer-core"
)

// NewErrorRenderertemplate returns a renderer which generates items with passed error message.
func NewErrorRenderer(template core.Template, nodeId string, err error) core.Renderer {
	return &ErrorRenderer{
		template: template,
		nodeId:   nodeId,
		err:      err,
	}
}

// Size returns size of entire error message box.
func (renderer *ErrorRenderer) Size() core.Size {
	return core.Size{
		Height: 124,
		Width:  400,
	}
}

// Content returns errpr message passed at initialization as items.
// Together with a title and a timestamp.
func (renderer *ErrorRenderer) Content() (string, error) {
	data := errorData{
		RenderId:  utils.NewId(),
		NodeId:    renderer.nodeId,
		Message:   renderer.err.Error(),
		TimeStamp: time.Now().Format("2006-01-02 15:04:05 MST"),
	}
	return renderer.template.RenderWith(data)
}

// ObserveDataSource has no effect for error renderer, because there's no datasource.
func (renderer *ErrorRenderer) ObserveDataSource(ctx context.Context) {

}

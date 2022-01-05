package syncsign

import (
	"errors"

	utils "github.com/tommzn/go-utils"
	core "github.com/tommzn/hdb-renderer-core"
)

// NewResponseRenderer returns a new renderer for eInk main content.
// Passed data have to contain all items which shloud be displayed.
func NewResponseRenderer(template core.Template, nodeId string, itemRenderer []core.Renderer) core.Renderer {
	return &ResponseRenderer{
		template:     template,
		nodeId:       nodeId,
		itemRenderer: itemRenderer,
	}
}

// Content returns the main layout for eInk display which includes
// renderer/node id and all passed items.
func (renderer *ResponseRenderer) Content() (string, error) {

	data := responseData{
		RenderId: utils.NewId(),
		NodeId:   renderer.nodeId,
		Items:    "",
	}

	items, err := renderer.contentFromItemRenderer()
	if err != nil {
		return "", err
	}
	if items == nil || *items == "" {
		return "", errors.New("No items has been rendered!")
	}
	data.Items = *items
	return renderer.template.RenderWith(data)
}

// ContentFromItemRenderer loops above all existing item renderers and returns a list of all generated items.
func (renderer *ResponseRenderer) contentFromItemRenderer() (*string, error) {

	content := ""
	errorStack := utils.NewErrorStack()
	for _, itemRenderer := range renderer.itemRenderer {
		items, err := itemRenderer.Content()
		errorStack.Append(err)
		content = appendItems(content, items)
	}
	return &content, errorStack.AsError()
}

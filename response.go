package syncsign

import (
	"context"
	"strings"

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

// Size is equal to size of 7.5 inch screen.
func (renderer *ResponseRenderer) Size() core.Size {
	return core.Size{
		Height: 528,
		Width:  880,
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
	data.Items = *items
	return renderer.template.RenderWith(data)
}

// ObserveDataSource has no effect for response renderer, because there's no datasource.
func (renderer *ResponseRenderer) ObserveDataSource(ctx context.Context) {

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

// AppendItems appends passed items to given content, separated by default JSON element separattor: ",".
// Leasing separators in content, or trailing separators in items will be removed.
func appendItems(content, items string) string {
	if items != "" {
		items = strings.TrimPrefix(items, ",")
		items = strings.TrimSuffix(items, ",")
		content = content + "," + items
	}
	return strings.TrimPrefix(strings.TrimSuffix(content, ","), ",")
}

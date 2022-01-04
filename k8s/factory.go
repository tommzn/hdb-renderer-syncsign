package main

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-renderer-core"
	syncsign "github.com/tommzn/hdb-renderer-syncsign"
)

func newFactory(conf config.Config, logger log.Logger) *factory {
	return &factory{
		conf:             conf,
		logger:           logger,
		responseRenderer: make(map[string]core.Renderer),
	}
}

func (f *factory) newResponseRendererTemplate() core.Template {
	if f.responseTemplate == nil {
		f.responseTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.response.template")
	}
	return f.responseTemplate
}

func (f *factory) newErrorRendererTemplate() core.Template {
	if f.errorTemplate == nil {
		f.errorTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.error_response.template")
	}
	return f.errorTemplate
}

func (f *factory) newIndoorClimateTemplate() core.Template {
	if f.indoorClimateTemplate == nil {
		f.indoorClimateTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.indoorclimate.template")
	}
	return f.indoorClimateTemplate
}

func (f *factory) newErrorRenderer(nodeId string, err error) core.Renderer {
	return syncsign.NewErrorRenderer(f.newErrorRendererTemplate(), nodeId, err)
}

func (f *factory) newResponseRenderer(nodeId string) core.Renderer {
	if _, ok := f.responseRenderer[nodeId]; !ok {
		f.responseRenderer[nodeId] = syncsign.NewResponseRenderer(f.newResponseRendererTemplate(), nodeId, f.itemRenderer())
	}
	return f.responseRenderer[nodeId]
}

func (f *factory) newIndoorClimateRenderer() core.Renderer {
	if f.indoorClimateRenderer == nil {
		f.indoorClimateRenderer = syncsign.NewIndoorClimateRenderer(f.conf, f.logger, f.newIndoorClimateTemplate(), f.newIndoorClimateDataSource())
	}
	return f.indoorClimateRenderer
}

func (f *factory) newIndoorClimateDataSource() core.DataSource {
	if f.indoorClimateDataSource == nil {
		f.indoorClimateDataSource = newDataSourceMock(f.indoorClimateDevices())
		f.indoorClimateDataSource.(*dataSourceMock).initMessages()
	}
	return f.indoorClimateDataSource
}

func (f *factory) indoorClimateDevices() []string {

	devices := []string{}
	devicesCfg := f.conf.GetAsSliceOfMaps("hdb.indoorclimate.devices")
	for _, device := range devicesCfg {
		if id, ok := device["id"]; ok {
			devices = append(devices, id)
		}
	}
	return devices
}

func (f *factory) itemRenderer() []core.Renderer {
	if f.indoorClimateRenderer == nil {
		f.newIndoorClimateRenderer()
	}
	return []core.Renderer{f.indoorClimateRenderer}
}

func (f *factory) newDisplayConfig() *syncsign.DisplayConfig {
	if f.displayConfig == nil {
		f.displayConfig = syncsign.NewDisplayConfig(f.conf)
	}
	return f.displayConfig
}

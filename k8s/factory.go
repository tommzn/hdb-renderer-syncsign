package main

import (
	"context"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-renderer-core"
	syncsign "github.com/tommzn/hdb-renderer-syncsign"
)

func newFactory(conf config.Config, logger log.Logger, ctx context.Context) *factory {
	return &factory{
		conf:             conf,
		logger:           logger,
		ctx:              ctx,
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
		f.errorTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.error.template")
	}
	return f.errorTemplate
}

func (f *factory) newIndoorClimateTemplate() core.Template {
	if f.indoorClimateTemplate == nil {
		f.indoorClimateTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.indoorclimate.template")
	}
	return f.indoorClimateTemplate
}

func (f *factory) newTimestampTemplate() core.Template {
	if f.timestampTemplate == nil {
		f.timestampTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.timestamp.template")
	}
	return f.timestampTemplate
}

func (f *factory) newTimestampRenderer() core.Renderer {
	return syncsign.NewTimestampRenderer(f.newTimestampTemplate())
}

func (f *factory) newErrorRenderer(err error) core.Renderer {
	return syncsign.NewErrorRenderer(f.newErrorRendererTemplate(), err)
}

func (f *factory) newErrorResponseRenderer(nodeId string, err error) core.Renderer {
	itemRenderer := []core.Renderer{
		f.newErrorRenderer(err),
		f.newTimestampRenderer(),
	}
	return syncsign.NewResponseRenderer(f.newResponseRendererTemplate(), nodeId, itemRenderer)
}

func (f *factory) newResponseRenderer(nodeId string) core.Renderer {
	if _, ok := f.responseRenderer[nodeId]; !ok {
		itemRenderer := []core.Renderer{
			f.newIndoorClimateRenderer(),
			f.newTimestampRenderer(),
		}
		f.responseRenderer[nodeId] = syncsign.NewResponseRenderer(f.newResponseRendererTemplate(), nodeId, itemRenderer)
	}
	return f.responseRenderer[nodeId]
}

func (f *factory) newIndoorClimateRenderer() core.Renderer {
	if f.indoorClimateRenderer == nil {
		renderer := syncsign.NewIndoorClimateRenderer(f.conf, f.logger, f.newIndoorClimateTemplate(), f.newIndoorClimateDataSource())
		go renderer.ObserveDataSource(f.ctx)
		f.indoorClimateRenderer = renderer
	}
	return f.indoorClimateRenderer
}

func (f *factory) newIndoorClimateDataSource() core.DataSource {
	if f.indoorClimateDataSource == nil {
		f.indoorClimateDataSource = newDataSourceMock(f.indoorClimateDevices())
		f.indoorClimateDataSource.(*dataSourceMock).initMessages()
		go f.indoorClimateDataSource.(*dataSourceMock).Run(f.ctx)
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

func (f *factory) newDisplayConfig() *syncsign.DisplayConfig {
	if f.displayConfig == nil {
		f.displayConfig = syncsign.NewDisplayConfig(f.conf)
	}
	return f.displayConfig
}

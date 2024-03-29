package main

import (
	"context"
	"sync"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
	datasource "github.com/tommzn/hdb-message-client"
	core "github.com/tommzn/hdb-renderer-core"
	syncsign "github.com/tommzn/hdb-renderer-syncsign"
)

func newFactory(conf config.Config, logger log.Logger, ctx context.Context) *factory {
	return &factory{
		conf:             conf,
		logger:           logger,
		ctx:              ctx,
		wg:               &sync.WaitGroup{},
		responseRenderer: make(map[string]core.Renderer),
		datasources:      []datasource.Client{},
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

func (f *factory) newBillingReportTemplate() core.Template {
	if f.billingReportTemplate == nil {
		f.billingReportTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.billingreport.template")
	}
	return f.billingReportTemplate
}

func (f *factory) newCurrentWeatherTemplate() core.Template {
	if f.currentWeatherTemplate == nil {
		f.currentWeatherTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.weather.template.current")
	}
	return f.currentWeatherTemplate
}

func (f *factory) newForeCastWeatherTemplate() core.Template {
	if f.forecastWeatherTemplate == nil {
		f.forecastWeatherTemplate = core.NewFileTemplateFromConfig(f.conf, "hdb.template_dir", "hdb.weather.template.forecast")
	}
	return f.forecastWeatherTemplate
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
			f.newBillingReportRenderer(),
			f.newWeatherRenderer(),
			f.newTimestampRenderer(),
		}
		f.responseRenderer[nodeId] = syncsign.NewResponseRenderer(f.newResponseRendererTemplate(), nodeId, itemRenderer)
	}
	return f.responseRenderer[nodeId]
}

func (f *factory) newIndoorClimateRenderer() core.Renderer {
	if f.indoorClimateRenderer == nil {
		renderer := syncsign.NewIndoorClimateRenderer(f.conf, f.logger, f.newIndoorClimateTemplate(), f.newDataSource())
		go renderer.ObserveDataSource(f.ctx)
		f.indoorClimateRenderer = renderer
	}
	return f.indoorClimateRenderer
}

func (f *factory) newBillingReportRenderer() core.Renderer {
	if f.billingReportRenderer == nil {
		renderer := syncsign.NewBillingReportRenderer(f.conf, f.logger, f.newBillingReportTemplate(), f.newDataSource())
		go renderer.ObserveDataSource(f.ctx)
		f.billingReportRenderer = renderer
	}
	return f.billingReportRenderer
}

func (f *factory) newWeatherRenderer() core.Renderer {
	if f.weatherRenderer == nil {
		renderer := syncsign.NewWeatherRenderer(f.conf, f.logger, f.newCurrentWeatherTemplate(), f.newForeCastWeatherTemplate(), f.newDataSource())
		go renderer.ObserveDataSource(f.ctx)
		f.weatherRenderer = renderer
	}
	return f.weatherRenderer
}

func (f *factory) newDataSource() core.DataSource {
	dataSource := datasource.New(f.conf, f.logger)
	f.wg.Add(1)
	go dataSource.Run(f.ctx, f.wg)
	f.datasources = append(f.datasources, dataSource)
	return dataSource
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

func (f *factory) dataSourceMetrics() map[int][]metrics.Measurement {

	dataSourceMetrics := make(map[int][]metrics.Measurement)
	for id, datasource := range f.datasources {
		dataSourceMetrics[id] = datasource.Metrics()
	}
	return dataSourceMetrics
}

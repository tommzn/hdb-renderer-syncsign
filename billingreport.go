package syncsign

import (
	"context"
	"errors"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	dsclient "github.com/tommzn/hdb-message-client"
	core "github.com/tommzn/hdb-renderer-core"
)

// NewBillingReportRenderer returns a renderer which generates items for AWS billing reports.
func NewBillingReportRenderer(conf config.Config, logger log.Logger, template core.Template, datasource core.DataSource) *BillingReportRenderer {

	anchor := anchorFromConfig(conf, "hdb.billingreport.anchor")
	reportCurrency := conf.Get("hdb.billingreport.report_currency", config.AsStringPtr("USD"))
	displayCurrency := conf.Get("hdb.billingreport.display_currency", config.AsStringPtr("USD"))
	return &BillingReportRenderer{
		template:        template,
		anchor:          anchor,
		logger:          logger,
		reportCurrency:  *reportCurrency,
		displayCurrency: *displayCurrency,
		datasource:      datasource,
		exchangeRates:   make(map[string]*events.ExchangeRate),
	}
}

// Content generates items for billing reports received from used datasource.
func (renderer *BillingReportRenderer) Content() (string, error) {

	renderer.logDatasource()

	if renderer.billingReport == nil {
		if err := renderer.fetchEvents(); err != nil {
			return "", errors.New("No billing report available.")
		}
	}
	return renderer.template.RenderWith(renderer.billingReport)
}

func (renderer *BillingReportRenderer) logDatasource() {
	if ds, ok := renderer.datasource.(*dsclient.MessageClient); ok {
		renderer.logger.Debug(ds.String())
	} else {
		renderer.logger.Debug("Unable to log datasource.")
	}

}

// FetchEvents will retrieve latest billing report from used datasource and process it
// and will retrieve all avaiable exchange rates as well, if they're necessary.
func (renderer *BillingReportRenderer) fetchEvents() error {

	if renderer.reportCurrency != renderer.displayCurrency {
		exchangeRates, err := renderer.datasource.All(hdbcore.DATASOURCE_EXCHANGERATE)
		if err == nil {
			for _, exchangeRate := range exchangeRates {
				renderer.processEvent(exchangeRate)
			}
		}
	}

	billingReport, err := renderer.datasource.Latest(hdbcore.DATASOURCE_BILLINGREPORT)
	if err == nil {
		renderer.processEvent(billingReport)
	}
	return err
}

// ObserveDataSource will listen for new billing reports and exchange rate events, if report and display currency differs.
func (renderer *BillingReportRenderer) ObserveDataSource(ctx context.Context) {

	defer renderer.logger.Flush()

	filter := []hdbcore.DataSource{hdbcore.DATASOURCE_BILLINGREPORT}
	if renderer.reportCurrency != renderer.displayCurrency {
		filter = append(filter, hdbcore.DATASOURCE_EXCHANGERATE)
	}
	renderer.dataSourceChan = renderer.datasource.Observe(&filter)
	for {
		select {
		case message, ok := <-renderer.dataSourceChan:
			if !ok {
				renderer.logger.Error("Error at reading datasource channel. Stop observing!")
				return
			}
			renderer.logger.Debug("Event received from datasource.")
			renderer.logger.Debug(message)
			renderer.processEvent(message)
		case <-ctx.Done():
			renderer.logger.Info("Camceled, stop observing.")
			return
		}
	}
}

// ProcessEvent will store latest billing report and exchange rates for comtemt remdering.
func (renderer *BillingReportRenderer) processEvent(message proto.Message) {

	if billingReport, ok := message.(*events.BillingReport); ok {
		renderer.logger.Debugf("Receive new billing report for %s", billingReport.BillingPeriod)
		renderer.calculateBillingReport(billingReport)
	}
	if exchangeRates, ok := message.(*events.ExchangeRates); ok {
		renderer.assignExchangeRates(exchangeRates)
	}
}

// AssignExchangeRates will save passed exchange rate if it's relevant for billing report calculations.
func (renderer *BillingReportRenderer) assignExchangeRates(exchangeRates *events.ExchangeRates) {

	for _, exchangeRate := range exchangeRates.Rates {
		renderer.logger.Debugf("Receive new exchange rate %s/%s", exchangeRate.FromCurrency, exchangeRate.ToCurrency)
		renderer.assignExchangeRate(exchangeRate)
	}
}

// AssignExchangeRate will save passed exchange rate locally if there's no exchangre assigned yet or
// if passed exchange rate is newer.
func (renderer *BillingReportRenderer) assignExchangeRate(exchangeRate *events.ExchangeRate) {

	exchangeRateKey := keyForExchangeRate(renderer.reportCurrency, renderer.displayCurrency)
	if exchangeRateKey != keyForExchangeRate(exchangeRate.FromCurrency, exchangeRate.ToCurrency) {
		return
	}
	if currentExchangeRate, ok := renderer.exchangeRates[exchangeRateKey]; ok &&
		exchangeRate.Timestamp.AsTime().Before(currentExchangeRate.Timestamp.AsTime()) {
		return
	}
	renderer.exchangeRates[exchangeRateKey] = exchangeRate
}

// CalculateBillingReport summarizes net and tax amount to a total billing report amount.
func (renderer *BillingReportRenderer) calculateBillingReport(billingReport *events.BillingReport) {

	totalAmount := billingReportAmount{
		Amount:   0.0,
		Currency: renderer.reportCurrency,
	}
	for _, amount := range billingReport.BillingAmount {
		totalAmount.Amount += amount
	}
	for _, amount := range billingReport.TaxAmount {
		totalAmount.Amount += amount
	}
	renderer.billingReport = &billingReportData{
		Anchor: renderer.anchor,
		Period: billingReport.BillingPeriod,
		Amount: formatForCurrency(renderer.convertAmount(totalAmount)),
	}
}

// ConvertAmount will convert billing report currency into display currency if both differs
// and if am exchange rate for both currencies is available.
func (renderer *BillingReportRenderer) convertAmount(amount billingReportAmount) billingReportAmount {
	if renderer.reportCurrency != renderer.displayCurrency {
		if exchangeRate, ok := renderer.exchangeRates[keyForExchangeRate(renderer.reportCurrency, renderer.displayCurrency)]; ok {
			return billingReportAmount{
				Amount:   amount.Amount * exchangeRate.Rate,
				Currency: renderer.displayCurrency,
			}
		}
	}
	return amount
}

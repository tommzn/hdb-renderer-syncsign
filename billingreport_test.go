package syncsign

import (
	"context"
	"github.com/stretchr/testify/suite"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type BillingReportTestSuite struct {
	suite.Suite
}

func TestBillingReportTestSuite(t *testing.T) {
	suite.Run(t, new(BillingReportTestSuite))
}

func (suite *BillingReportTestSuite) TestGenerateContent() {

	renderer := billingReportRendererForTest("fixtures/testconfig04.yml")

	content, err := renderer.Content()
	suite.Nil(err)
	assertTemplateHash(suite.Assert(), content, "ee81fe0a5b91648eab2b8398ad385e44e683d2f6")
}

func (suite *BillingReportTestSuite) TestWithoutCurrencyConversion() {

	renderer := billingReportRendererForTest("fixtures/testconfig05.yml")

	content, err := renderer.Content()
	suite.Nil(err)
	assertTemplateHash(suite.Assert(), content, "d0bcc54397e73be3e89b0434ebde08838f6e4921")
}

func (suite *BillingReportTestSuite) TestGenerateContentByObservingDataSource() {

	renderer := billingReportRendererForTest("fixtures/testconfig04.yml")
	ctx, cancelFunc := context.WithCancel(context.Background())
	endChan := make(chan bool, 1)
	go func() {
		renderer.ObserveDataSource(ctx)
		endChan <- true
	}()
	time.Sleep(1 * time.Second)
	suite.NotNil(renderer.billingReport)
	suite.True(len(renderer.exchangeRates) > 0)

	content, err := renderer.Content()
	suite.Nil(err)
	assertTemplateHash(suite.Assert(), content, "ee81fe0a5b91648eab2b8398ad385e44e683d2f6")

	cancelFunc()
	select {
	case ok := <-endChan:
		suite.True(ok)
	case <-time.After(200 * time.Millisecond):
		suite.T().Error("DataSource observing doesn't end as expected!")
	}
}

func (suite *BillingReportTestSuite) TestStopDataSourceObservingOnClosedChannel() {

	renderer := billingReportRendererForTest("fixtures/testconfig04.yml")
	ctx, _ := context.WithCancel(context.Background())

	endChan := make(chan bool, 1)
	go func() {
		renderer.ObserveDataSource(ctx)
		endChan <- true
	}()
	time.Sleep(100 * time.Millisecond)

	close(renderer.datasource.(*datasourceMock).eventChan)
	select {
	case ok := <-endChan:
		suite.True(ok)
	case <-time.After(200 * time.Millisecond):
		suite.T().Error("DataSource observing doesn't end as expected!")
	}
}

func (suite *BillingReportTestSuite) TestAssignExchangeRates() {

	renderer := billingReportRendererForTest("fixtures/testconfig04.yml")
	fromCurrency := "USD"
	toCurrency := "EUR"
	exchangeRateKey := keyForExchangeRate(fromCurrency, toCurrency)
	exchangeRate1 := &events.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         0.8345,
		Timestamp:    timestamppb.New(time.Now()),
	}
	exchangeRate2 := &events.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         1.2325,
		Timestamp:    timestamppb.New(time.Now().Add(-100 * time.Second)),
	}

	renderer.assignExchangeRate(exchangeRate1)
	assignedRate1, ok1 := renderer.exchangeRates[exchangeRateKey]
	suite.True(ok1)
	suite.Equal(exchangeRate1.Rate, assignedRate1.Rate)

	renderer.assignExchangeRate(exchangeRate2)
	assignedRate2, ok2 := renderer.exchangeRates[exchangeRateKey]
	suite.True(ok2)
	suite.Equal(exchangeRate1.Rate, assignedRate2.Rate)
}

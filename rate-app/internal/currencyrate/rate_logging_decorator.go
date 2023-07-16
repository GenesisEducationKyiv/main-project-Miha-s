package currencyrate

import (
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"net/http"
)

type RateLoggingDecorator struct {
	executor NamedHttpRateExecutor
}

func NewRateLoggingDecorator(executor NamedHttpRateExecutor) *RateLoggingDecorator {
	return &RateLoggingDecorator{
		executor: executor,
	}
}

func (decorator *RateLoggingDecorator) GenerateHttpRequest(currency *models.Currency) (*http.Request, error) {
	return decorator.executor.GenerateHttpRequest(currency)
}

func (decorator *RateLoggingDecorator) ExtractRate(resp *http.Response, currency *models.Currency) (models.Rate, error) {
	rate, err := decorator.executor.ExtractRate(resp, currency)
	logger.Log.Infof("Rate from %v is %v", decorator.executor.Name(), rate.Value)
	return rate, err
}

func (decorator *RateLoggingDecorator) Name() string {
	return decorator.executor.Name()
}

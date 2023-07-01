package rateProvidersTest

import (
	"btc-test-task/internal/helpers/models"
	"btc-test-task/internal/rateProviders"
)

type RateProviderStub struct {
	RateError error
	Rate      models.Rate
}

func (provider *RateProviderStub) GetCurrentRate(currencyFrom string, currencyTo string) (models.Rate, error) {
	return provider.Rate, provider.RateError
}

func (provider *RateProviderStub) SetNext(nextProvider rateProviders.RateProvider) {

}

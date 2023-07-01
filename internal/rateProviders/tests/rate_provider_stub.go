package rateProvidersTest

import (
	"btc-test-task/internal/helpers/models"
)

type RateProviderStub struct {
	RateError error
	Rate      models.Rate
}

func (provider *RateProviderStub) GetCurrentRate(currency *models.Currency) (models.Rate, error) {
	return provider.Rate, provider.RateError
}

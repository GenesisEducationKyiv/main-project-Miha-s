package rateProvidersTest

import (
	models2 "btc-test-task/internal/common/models"
)

type RateProviderStub struct {
	RateError error
	Rate      models2.Rate
}

func (provider *RateProviderStub) GetCurrentRate(currency *models2.Currency) (models2.Rate, error) {
	return provider.Rate, provider.RateError
}

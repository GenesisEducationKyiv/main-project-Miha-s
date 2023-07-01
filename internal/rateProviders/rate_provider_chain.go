package rateProviders

import "btc-test-task/internal/helpers/models"

type RateProviderChain interface {
	GetCurrentRate(currencyFrom string, currencyTo string) (models.Rate, error)
	SetNext(provider RateProviderChain)
}

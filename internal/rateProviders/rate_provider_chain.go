package rateProviders

import "btc-test-task/internal/helpers/models"

type RateProviderChain interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
	SetNext(provider RateProviderChain)
}

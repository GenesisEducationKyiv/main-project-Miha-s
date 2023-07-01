package handlers

import (
	"btc-test-task/internal/helpers/models"
)

type RateProvider interface {
	GetCurrentRate(currencyFrom string, currencyTo string) (models.Rate, error)
}

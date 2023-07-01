package handlers

import (
	"btc-test-task/internal/helpers/models"
)

type RateProvider interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
}

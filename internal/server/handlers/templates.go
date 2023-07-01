package handlers

import (
	"btc-test-task/internal/helpers/models"
)

type Templates interface {
	CurrencyRate(rate models.Rate) string
}

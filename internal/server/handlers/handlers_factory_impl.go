package handlers

import (
	"btc-test-task/internal/common/configuration/config"
)

type HandlersFactoryImpl struct {
	services     Services
	currencyFrom string
	currencyTo   string
}

func NewHandlersFactoryImpl(conf *config.Config, services Services) *HandlersFactoryImpl {
	return &HandlersFactoryImpl{
		services:     services,
		currencyFrom: conf.CurrencyFrom,
		currencyTo:   conf.CurrencyTo,
	}
}

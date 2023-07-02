package handlers

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/types"
)

type HandlersFactoryImpl struct {
	services     *types.Services
	currencyFrom string
	currencyTo   string
}

func NewHandlersFactoryImpl(conf *config.Config, services *types.Services) (*HandlersFactoryImpl, error) {
	newHandlersFactoryImpl := new(HandlersFactoryImpl)
	err := newHandlersFactoryImpl.init(conf, services)
	if err != nil {
		return nil, err
	}
	return newHandlersFactoryImpl, nil
}

func (factory *HandlersFactoryImpl) init(conf *config.Config, services *types.Services) error {
	factory.services = services
	factory.currencyFrom = conf.CurrencyFrom
	factory.currencyTo = conf.CurrencyTo
	return nil
}

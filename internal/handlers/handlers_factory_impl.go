package handlers

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/types"
)

type HandlersFactoryImpl struct {
	services *types.Services
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
	return nil
}

package handlers

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/types"
)

type HandlersFactoryImpl struct {
	services *types.Services
}

func (factory *HandlersFactoryImpl) Init(conf *config.Config, services *types.Services) error {
	factory.services = services
	return nil
}

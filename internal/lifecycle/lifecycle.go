package lifecycle

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/handlers"
	"btc-test-task/internal/logger"
	"btc-test-task/internal/server"
	"btc-test-task/internal/types"

	"os"
	"os/signal"
	"syscall"
)

type Lifecycle struct {
	services         types.Services
	handlers_factory handlers.HandlersFactory
	server           server.Server
	config           config.Config
}

func (lifecycle *Lifecycle) Init(conf *config.Config) error {
	lifecycle.config = *conf
	logger.Init(conf)
	lifecycle.services.EmailSender.Init(conf)
	lifecycle.services.CourseAccessor.Init(conf)
	lifecycle.services.EmailStorage.Init(conf)
	lifecycle.handlers_factory = &handlers.HandlersFactoryImpl{}
	lifecycle.handlers_factory.Init(conf, &lifecycle.services)

	lifecycle.server.Init(conf, lifecycle.handlers_factory)
	return nil
}

func (lifecycle *Lifecycle) Run() error {
	done := make(chan error, 1)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	defer lifecycle.services.EmailStorage.Close()
	go lifecycle.server.Run()

	select {
	case <-signals:
	case <-done:
	}
	return nil
}

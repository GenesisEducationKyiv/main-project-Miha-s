package lifecycle

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/email_sender"
	"btc-test-task/internal/emails_storage"
	"btc-test-task/internal/handlers"
	"btc-test-task/internal/logger"
	"btc-test-task/internal/rate_accessors"
	"btc-test-task/internal/server"
	"btc-test-task/internal/templates"
	"btc-test-task/internal/types"
	"fmt"
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

	lifecycle.services.Templates = &templates.TemplatesImpl{}
	err := lifecycle.services.Templates.Init(conf)
	if err != nil {
		return err
	}
	lifecycle.services.EmailSender = &email_sender.EmailSenderImpl{}
	err = lifecycle.services.EmailSender.Init(conf)
	if err != nil {
		return err
	}

	lifecycle.services.RateAccessor = &rate_accessors.CoinApI{}
	err = lifecycle.services.RateAccessor.Init(conf)
	if err != nil {
		return err
	}

	lifecycle.services.EmailStorage = &emails_storage.EmailsStorageImpl{}
	err = lifecycle.services.EmailStorage.Init(conf)
	if err != nil {
		return err
	}

	lifecycle.handlers_factory = &handlers.HandlersFactoryImpl{}
	err = lifecycle.handlers_factory.Init(conf, &lifecycle.services)
	if err != nil {
		return err
	}

	err = lifecycle.server.Init(conf, lifecycle.handlers_factory)
	if err != nil {
		return err
	}

	return nil
}

func (lifecycle *Lifecycle) Run() error {
	done := make(chan error, 1)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	defer lifecycle.services.EmailStorage.Close()

	go func() {
		done <- lifecycle.server.Run()
	}()

	select {
	case <-signals:
		logger.LogInfo("Signal was received, shutting down...")
	case err := <-done:
		if err != nil {
			logger.LogErrorStr(fmt.Sprintf("Server crashed with error %v", err))
		} else {
			logger.LogInfo("Server finished its work, shutting down...")
		}
	}
	return nil
}

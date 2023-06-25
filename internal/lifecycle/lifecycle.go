package lifecycle

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/emailSender"
	"btc-test-task/internal/emailsStorage"
	"btc-test-task/internal/handlers"
	"btc-test-task/internal/logger"
	"btc-test-task/internal/rateAccessors"
	"btc-test-task/internal/server"
	"btc-test-task/internal/templates"
	"btc-test-task/internal/types"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Lifecycle struct {
	services        types.Services
	handlersFactory handlers.HandlersFactory
	server          *server.Server
	config          config.Config
}

func (lifecycle *Lifecycle) Init(conf *config.Config) error {
	lifecycle.config = *conf
	logger.Init(conf)
	err := error(nil)
	lifecycle.services.Templates, err = templates.NewSimpleTextTemplates(conf)
	if err != nil {
		return err
	}
	lifecycle.services.EmailSender, err = emailSender.NewGoMailSender(conf)
	if err != nil {
		return err
	}

	lifecycle.services.RateAccessor, err = rateAccessors.NewCoinAPI(conf)
	if err != nil {
		return err
	}

	lifecycle.services.EmailStorage, err = emailsStorage.NewJsonEmailsStorage(conf)
	if err != nil {
		return err
	}

	lifecycle.handlersFactory, err = handlers.NewHandlersFactoryImpl(conf, &lifecycle.services)
	if err != nil {
		return err
	}

	lifecycle.server, err = server.NewServer(conf, lifecycle.handlersFactory)
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

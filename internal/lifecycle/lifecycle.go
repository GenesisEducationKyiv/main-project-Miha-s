package lifecycle

import (
	"btc-test-task/internal/emailSender"
	"btc-test-task/internal/emailsRepository"
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/templates"
	"btc-test-task/internal/helpers/types"
	"btc-test-task/internal/helpers/validators"
	"btc-test-task/internal/rateProviders"
	"btc-test-task/internal/server"
	"btc-test-task/internal/server/handlers"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
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
		return errors.Wrap(err, "Init")
	}
	lifecycle.services.EmailSender, err = emailSender.NewGoMailSender(conf)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	CoinGeckoRateProvider := rateProviders.NewHttpRateProvider(rateProviders.NewCoinGeckoExecutor(conf))
	CoinAPIRateProvider := rateProviders.NewHttpRateProvider(rateProviders.NewCoinAPIExecutor(conf))
	CoinGeckoRateProvider.SetNext(CoinAPIRateProvider)
	lifecycle.services.RateProvider = CoinGeckoRateProvider

	lifecycle.services.EmailsRepository, err = emailsRepository.NewJsonEmailsStorage(conf, new(validators.RegexEmailValidator))
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.handlersFactory, err = handlers.NewHandlersFactoryImpl(conf, &lifecycle.services)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.server, err = server.NewServer(conf, lifecycle.handlersFactory)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	logger.Log.Infof("The server is listening on port: %v", conf.Port)

	return nil
}

func (lifecycle *Lifecycle) Run() error {
	done := make(chan error, 1)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	defer lifecycle.services.EmailsRepository.Close()

	go func() {
		done <- lifecycle.server.Run()
	}()

	select {
	case <-signals:
		logger.Log.Info("Signal was received, shutting down...")
	case err := <-done:
		if err != nil {
			logger.Log.Errorf("Server crashed with error %v", err)
		} else {
			logger.Log.Info("Server finished its work, shutting down...")
		}
	}
	return nil
}

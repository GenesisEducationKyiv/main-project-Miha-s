package lifecycle

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/currencyrate"
	"btc-test-task/internal/email"
	"btc-test-task/internal/email/templates"
	"btc-test-task/internal/repository"
	"btc-test-task/internal/repository/validators"
	"btc-test-task/internal/server"
	"btc-test-task/internal/server/handlers"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

type Lifecycle struct {
	services        Services
	handlersFactory server.HandlersFactory
	server          *server.Server
	config          config.Config
}

func (lifecycle *Lifecycle) Init(conf *config.Config) error {
	lifecycle.config = *conf
	err := logger.Init(conf)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.services.Templates = templates.NewSimpleTextTemplates(conf)
	lifecycle.services.EmailSender = email.NewGoMailSender(conf)

	CoinGeckoRateProvider := currencyrate.NewHttpRateProvider(currencyrate.NewCoinGeckoExecutor(conf))
	CoinAPIRateProvider := currencyrate.NewHttpRateProvider(currencyrate.NewCoinAPIExecutor(conf))
	BinanceAPIrateProvider := currencyrate.NewHttpRateProvider(currencyrate.NewBinanceAPIExecutor(conf))
	CoinAPIRateProvider.SetNext(BinanceAPIrateProvider)
	CoinGeckoRateProvider.SetNext(CoinAPIRateProvider)
	lifecycle.services.RateProvider = CoinGeckoRateProvider

	lifecycle.services.EmailsRepository, err = repository.NewJsonEmailsStorage(conf, new(validators.RegexEmailValidator))
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.handlersFactory = handlers.NewHandlersFactoryImpl(conf, &lifecycle.services)

	lifecycle.server = server.NewServer(conf, lifecycle.handlersFactory)

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

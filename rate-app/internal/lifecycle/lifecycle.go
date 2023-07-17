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
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type Lifecycle struct {
	services         Services
	handlersFactory  server.HandlersFactory
	loggerWriter     logger.LoggerWriter
	middlewareLogger *server.RabbitMQLogger
	server           *server.Server
	config           config.Config
}

func (lifecycle *Lifecycle) Init(conf *config.Config) error {
	lifecycle.config = *conf
	var err error

	lifecycle.loggerWriter, err = composeLoggers(conf)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	err = logger.Init(conf, lifecycle.loggerWriter)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.services.Templates = templates.NewSimpleTextTemplates(conf)
	lifecycle.services.EmailSender = email.NewGoMailSender(conf)

	lifecycle.services.RateProvider, err = composeRateProvider(conf)
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.services.EmailsRepository, err = repository.NewJsonEmailsStorage(conf, new(validators.RegexEmailValidator))
	if err != nil {
		return errors.Wrap(err, "Init")
	}

	lifecycle.handlersFactory = handlers.NewHandlersFactoryImpl(conf, &lifecycle.services)
	middlewareLogger, err := server.NewRabbitMQLogger(conf)
	if err != nil {
		return errors.Wrap(err, "Init")
	}
	lifecycle.middlewareLogger = middlewareLogger
	lifecycle.server = server.NewServer(conf, lifecycle.handlersFactory, middlewareLogger.MiddlewareLogger())

	logger.Log.Infof("The server is listening on port: %v", conf.Port)

	return nil
}

func composeLoggers(conf *config.Config) (logger.LoggerWriter, error) {
	consoleLogger := logger.NewLoggerWriterChain(logger.NewConsoleLoggerWriter())

	if len(conf.LogFile) != 0 {
		fileLogger, err := logger.NewFileLoggerWriter(conf)
		if err != nil {
			return nil, err
		}
		consoleLogger.SetNext(logger.NewLoggerWriterChain(fileLogger))
	}

	return consoleLogger, nil
}

func composeRateProvider(conf *config.Config) (handlers.RateProvider, error) {
	CoinGeckoRateProvider := currencyrate.NewHttpRateProvider(
		currencyrate.NewRateLoggingDecorator(currencyrate.NewCoinGeckoExecutor(conf)), http.DefaultClient)
	CoinAPIRateProvider := currencyrate.NewHttpRateProvider(
		currencyrate.NewRateLoggingDecorator(currencyrate.NewCoinAPIExecutor(conf)), http.DefaultClient)
	BinanceAPIrateProvider := currencyrate.NewHttpRateProvider(
		currencyrate.NewRateLoggingDecorator(currencyrate.NewBinanceAPIExecutor(conf)), http.DefaultClient)
	CoinAPIRateProvider.SetNext(BinanceAPIrateProvider)
	CoinGeckoRateProvider.SetNext(CoinAPIRateProvider)

	rateCache, err := currencyrate.NewRateCache(conf, CoinGeckoRateProvider,
		cache.New(conf.RateCacheDuration, conf.RateCacheDuration))
	if err != nil {
		return nil, err
	}
	return rateCache, nil
}

func (lifecycle *Lifecycle) Run() error {
	done := make(chan error, 1)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	defer lifecycle.services.EmailsRepository.Close()
	defer lifecycle.loggerWriter.Close()
	defer lifecycle.middlewareLogger.Close()
	defer lifecycle.server.Shutdown()

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

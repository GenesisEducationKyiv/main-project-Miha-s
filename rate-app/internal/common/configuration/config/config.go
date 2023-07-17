package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    uint
	EmailStorageName        string
	EmailStoragePath        string
	EmailToSendFrom         string
	EmailToSendFromPassword string
	EmailServiceUrl         string
	EmailServicePort        int
	EmailSubject            string
	CoinGeckoAPIUrl         string
	CoinAPIUrl              string
	CoinAPIKey              string
	BinanceAPIUrl           string
	CurrencyFrom            string
	CurrencyTo              string
	RateCacheDuration       time.Duration
	LogLevel                string
	LogFile                 string
	RabbitUsername          string
	RabbitPassword          string
	RabbitUrlPort           string
}

func (conf *Config) LoadFromENV(envFilePath string) error {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return errors.Wrap(err, "LoadFromENV")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port < 0 {
		return errors.Wrap(err, "LoadFromENV")
	}
	conf.Port = uint(port)

	conf.CoinGeckoAPIUrl = os.Getenv("COINGECKO_URL")
	conf.CoinAPIUrl = os.Getenv("COINAPI_URL")
	conf.CoinAPIKey = os.Getenv("COINAPI_KEY")
	conf.BinanceAPIUrl = os.Getenv("BINANCEAPI_URL")
	conf.CurrencyFrom = os.Getenv("CURRENCY_FROM")
	conf.CurrencyTo = os.Getenv("CURRENCY_TO")
	duration := os.Getenv("RATE_CACHE_DURATION")
	conf.RateCacheDuration, err = time.ParseDuration(duration)
	if err != nil {
		return errors.Wrap(err, "LoadFromENV")
	}

	conf.EmailStorageName = os.Getenv("EMAIL_STORAGE_NAME")
	conf.EmailStoragePath = os.Getenv("EMAIL_STORAGE_PATH")
	conf.EmailSubject = os.Getenv("EMAIL_SUBJECT")

	conf.EmailServiceUrl = os.Getenv("EMAIL_SERVICE_URL")
	conf.EmailServicePort, _ = strconv.Atoi(os.Getenv("EMAIL_SERVICE_PORT"))
	conf.EmailToSendFrom = os.Getenv("EMAIL_TO_SEND_FROM")
	conf.EmailToSendFromPassword = os.Getenv("EMAIL_TO_SEND_FROM_PASSWORD")

	conf.LogLevel = os.Getenv("LOG_LEVEL")
	conf.LogFile = os.Getenv("LOG_FILE")

	conf.RabbitUsername = os.Getenv("RABBIT_USERNAME")
	conf.RabbitPassword = os.Getenv("RABBIT_PASSWORD")
	conf.RabbitUrlPort = os.Getenv("RABBIT_URL_PORT")

	return nil
}

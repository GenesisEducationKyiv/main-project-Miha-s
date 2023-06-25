package config

import (
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    uint
	EmailStoragePath        string
	EmailToSendFrom         string
	EmailToSendFromPassword string
	EmailServiceUrl         string
	EmailServicePort        int
	EmailSubject            string
	CoinAPIUrl              string
	CoinAPIKey              string
	CurrencyFrom            string
	CurrencyTo              string
	LogLevel                string
	LogFile                 string
}

func (conf *Config) LoadFromENV() error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.Wrap(err, "LoadFromENV: ")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return errors.Wrap(err, "LoadFromENV: ")
	}
	conf.Port = uint(port)

	conf.CoinAPIUrl = os.Getenv("COINAPI_URL")
	conf.CoinAPIKey = os.Getenv("COINAPI_KEY")
	conf.CurrencyFrom = os.Getenv("CURRENCY_FROM")
	conf.CurrencyTo = os.Getenv("CURRENCY_TO")

	conf.EmailStoragePath = os.Getenv("EMAIL_STORAGE_PATH")
	conf.EmailSubject = os.Getenv("EMAIL_SUBJECT")

	conf.EmailServiceUrl = os.Getenv("EMAIL_SERVICE_URL")
	conf.EmailServicePort, _ = strconv.Atoi(os.Getenv("EMAIL_SERVICE_PORT"))
	conf.EmailToSendFrom = os.Getenv("EMAIL_TO_SEND_FROM")
	conf.EmailToSendFromPassword = os.Getenv("EMAIL_TO_SEND_FROM_PASSWORD")

	conf.LogLevel = os.Getenv("LOG_LEVEL")
	conf.LogFile = os.Getenv("LOG_FILE")

	return nil
}
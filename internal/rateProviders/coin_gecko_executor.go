package rateProviders

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/customErrors"
	"btc-test-task/internal/helpers/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type CoinGeckoExecutor struct {
	url string
}

var currenciesCoinGecko = map[string]string{
	"BTC": "bitcoin",
	"UAH": "uah",
}

type CoinGeckoResponse map[string]map[string]float64

func NewCoinGeckoExecutor(conf *config.Config) *CoinGeckoExecutor {
	return &CoinGeckoExecutor{
		url: conf.CoinGeckoAPIUrl,
	}
}

func (api *CoinGeckoExecutor) generateEndpoint(currency *models.Currency) (string, error) {
	return fmt.Sprintf(api.url, currency.CurrencyFrom, currency.CurrencyTo), nil
}

func (api *CoinGeckoExecutor) GenerateHttpRequest(currency *models.Currency) (*http.Request, error) {
	endpoint, err := api.generateEndpoint(currency)
	if err != nil {
		return nil, errors.Wrap(err, "GenerateCurrencyRequest")
	}
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "GenerateCurrencyRequest")
	}
	return req, err
}

func (api *CoinGeckoExecutor) ExtractRate(resp []byte, currency *models.Currency) (models.Rate, error) {
	rate := models.Rate{}

	var responseData CoinGeckoResponse
	err := json.NewDecoder(bytes.NewReader(resp)).Decode(&responseData)
	if err != nil {
		return rate, errors.Wrap(customErrors.ErrFailedToGetRate, "extractRate")
	}

	rate.Value = responseData[currenciesCoinGecko[currency.CurrencyFrom]][currenciesCoinGecko[currency.CurrencyTo]]
	if rate.Value == 0 {
		return rate, errors.Wrap(customErrors.ErrFailedToGetRate, "extractRate")
	}
	return rate, nil
}

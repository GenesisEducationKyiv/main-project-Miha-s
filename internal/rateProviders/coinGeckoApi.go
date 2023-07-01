package rateProviders

import (
	"btc-test-task/internal/helpers/config"
	errors2 "btc-test-task/internal/helpers/errors"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type CoinGeckoApi struct {
	url          string
	nextProvider RateProviderChain
}

type CoinGeckoResponse map[string]map[string]float64

func NewCoinGeckoApi(conf *config.Config) (*CoinGeckoApi, error) {
	return &CoinGeckoApi{
		url: conf.CoinGeckoAPIUrl,
	}, nil
}

func (api *CoinGeckoApi) SetNext(provider RateProviderChain) {
	api.nextProvider = provider
}

func (api *CoinGeckoApi) extractRate(resp *http.Response, currencyFrom string, currencyTo string) (models.Rate, error) {
	rate := models.Rate{}

	var responseData CoinGeckoResponse
	err := json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return rate, errors.Wrap(errors2.ErrFailedToGetRate, "extractRate")
	}

	rate.Value = responseData[strings.ToLower(currencyFrom)][strings.ToLower(currencyTo)]

	return rate, nil
}

func (api *CoinGeckoApi) generateEndpoint(currencyFrom string, currencyTo string) (string, error) {
	return fmt.Sprintf(api.url, currencyFrom, currencyTo), nil
}

func (api *CoinGeckoApi) generateCurrencyRequest(currencyFrom string, currencyTo string) (*http.Request, error) {
	endpoint, err := api.generateEndpoint(currencyFrom, currencyTo)
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

func (api *CoinGeckoApi) getRate(currencyFrom string, currencyTo string) (models.Rate, error) {
	currentRate := models.Rate{}
	req, err := api.generateCurrencyRequest(currencyFrom, currencyTo)
	if err != nil {
		logger.Log.Error(err)
		return currentRate, errors.Wrap(err, "GetCurrentRate")
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return currentRate, errors.Wrap(errors2.ErrFailedToGetRate, "GetCurrentRate")
	}

	currentRate, err = api.extractRate(res, currencyFrom, currencyTo)
	if err != nil {
		return currentRate, errors.Wrap(err, "GetCurrentRate")
	}

	logger.Log.Infof("CoinGeckoAPI rate %v", currentRate.Value)
	return currentRate, nil
}

func (api *CoinGeckoApi) GetCurrentRate(currencyFrom string, currencyTo string) (models.Rate, error) {
	rate, err := api.getRate(currencyFrom, currencyTo)
	if err != nil && api.nextProvider != nil {
		rate, err = api.nextProvider.GetCurrentRate(currencyFrom, currencyTo)
		return rate, err
	}
	return rate, err
}

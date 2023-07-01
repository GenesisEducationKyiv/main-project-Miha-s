package rateProviders

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/models"
	"encoding/json"
	"fmt"

	"io"
	"net/http"

	"github.com/pkg/errors"
)

type CoinAPI struct {
	endpoint     string
	apiKey       string
	nextProvider RateProvider
}

func NewCoinAPI(conf *config.Config) (*CoinAPI, error) {
	return &CoinAPI{
		endpoint: conf.CoinAPIUrl,
		apiKey:   conf.CoinAPIKey,
	}, nil
}

func (api *CoinAPI) SetNext(provider RateProvider) {
	api.nextProvider = provider
}

func (api *CoinAPI) extractRate(resp *http.Response) (models.Rate, error) {
	rate := models.Rate{}
	jsonValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return rate, errors.Wrap(ErrFailedToGetRate, "extractRate")
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonValue, &dat); err != nil {
		return rate, errors.Wrap(ErrFailedToGetRate, "extractRate")
	}
	var ok bool
	rate.Value, ok = dat["rate"].(float64)
	if !ok {
		return rate, errors.Wrap(ErrFailedToGetRate, "extractRate")
	}
	return rate, nil
}

func (api *CoinAPI) generateEndpoint(currencyFrom string, currencyTo string) (string, error) {
	return fmt.Sprintf(api.endpoint, currencyFrom, currencyTo), nil
}

func (api *CoinAPI) generateCurrencyRequest(currencyFrom string, currencyTo string) (*http.Request, error) {
	endpoint, err := api.generateEndpoint(currencyFrom, currencyTo)
	if err != nil {
		return nil, errors.Wrap(err, "generateCurrencyRequest")
	}
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(ErrFailedToGetRate, "generateCurrencyRequest")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CoinAPI-Key", api.apiKey)

	return req, nil
}

func (api *CoinAPI) getRate(currencyFrom string, currencyTo string) (models.Rate, error) {
	currentRate := models.Rate{}
	req, err := api.generateCurrencyRequest(currencyFrom, currencyTo)
	if err != nil {
		logger.Log.Error(err)
		return currentRate, errors.Wrap(err, "GetCurrentRate")
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return currentRate, errors.Wrap(ErrFailedToGetRate, "GetCurrentRate")
	}

	currentRate, err = api.extractRate(res)
	if err != nil {
		return currentRate, errors.Wrap(err, "GetCurrentRate")
	}

	logger.Log.Infof("CoinAPI rate %v", currentRate.Value)
	return currentRate, nil
}

func (api *CoinAPI) GetCurrentRate(currencyFrom string, currencyTo string) (models.Rate, error) {
	rate, err := api.getRate(currencyFrom, currencyTo)
	if err != nil && api.nextProvider != nil {
		rate, err = api.nextProvider.GetCurrentRate(currencyFrom, currencyTo)
		return rate, err
	}
	return rate, err
}

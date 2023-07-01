package rateProviders

import (
	"btc-test-task/internal/helpers/customErrors"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/models"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type HttpRateExecutor interface {
	GenerateHttpRequest(currency *models.Currency) (*http.Request, error)
	ExtractRate(resp []byte, currency *models.Currency) (models.Rate, error)
}

type HttpRateProvider struct {
	nextProvider RateProviderChain
	executor     HttpRateExecutor
}

func NewHttpRateProvider(executor HttpRateExecutor) *HttpRateProvider {
	return &HttpRateProvider{
		executor: executor,
	}
}

func (api *HttpRateProvider) SetNext(provider RateProviderChain) {
	api.nextProvider = provider
}

func (api *HttpRateProvider) executeRateRequest(currency *models.Currency) (*http.Response, error) {
	req, err := api.executor.GenerateHttpRequest(currency)
	if err != nil {
		logger.Log.Error(err)
		return nil, errors.Wrap(err, "executeRateRequest")
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.Wrap(customErrors.ErrFailedToGetRate, "executeRateRequest")
	}
	return res, nil
}

func (api *HttpRateProvider) parseRateResponse(response *http.Response, currency *models.Currency) (models.Rate, error) {
	currentRate := models.Rate{}

	responseValue, err := io.ReadAll(response.Body)
	if err != nil {
		return currentRate, errors.Wrap(customErrors.ErrFailedToGetRate, "GetCurrentRate")
	}

	currentRate, err = api.executor.ExtractRate(responseValue, currency)
	if err != nil {
		logger.Log.Errorf("failed to extract rate, response: %v", string(responseValue))
		return currentRate, errors.Wrap(err, "GetCurrentRate")
	}

	return currentRate, nil
}

func (api *HttpRateProvider) getRate(currency *models.Currency) (models.Rate, error) {
	currentRate := models.Rate{}
	res, err := api.executeRateRequest(currency)
	defer res.Body.Close()
	if err != nil {
		return currentRate, err
	}

	currentRate, err = api.parseRateResponse(res, currency)
	if err != nil {
		return currentRate, err
	}

	logger.Log.Infof("CoinGeckoAPI rate %v", currentRate.Value)
	return currentRate, nil
}

func (api *HttpRateProvider) GetCurrentRate(currency *models.Currency) (models.Rate, error) {
	rate, err := api.getRate(currency)
	if err == nil {
		return rate, nil
	} else if api.nextProvider == nil {
		return rate, errors.Wrap(err, "GetCurrentRate")
	}

	rate, err = api.nextProvider.GetCurrentRate(currency)
	if err != nil {
		return rate, errors.Wrap(err, "GetCurrentRate")
	}
	return rate, nil

}

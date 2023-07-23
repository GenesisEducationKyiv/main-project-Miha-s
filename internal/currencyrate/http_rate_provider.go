package currencyrate

import (
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"net/http"

	"github.com/pkg/errors"
)

type NamedHttpRateExecutor interface {
	GenerateHttpRequest(currency *models.Currency) (*http.Request, error)
	ExtractRate(resp *http.Response, currency *models.Currency) (models.Rate, error)
	Name() string
}

type RateProviderChain interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
	SetNext(provider RateProviderChain)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpRateProvider struct {
	nextProvider RateProviderChain
	executor     NamedHttpRateExecutor
	client       HttpClient
}

func NewHttpRateProvider(executor NamedHttpRateExecutor, client HttpClient) *HttpRateProvider {
	return &HttpRateProvider{
		executor: executor,
		client:   client,
	}
}

func (api *HttpRateProvider) SetNext(provider RateProviderChain) {
	api.nextProvider = provider
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

func (api *HttpRateProvider) getRate(currency *models.Currency) (models.Rate, error) {
	currentRate := models.Rate{}
	res, err := api.executeRateRequest(currency)
	if err != nil {
		return currentRate, err
	}
	defer res.Body.Close()

	currentRate, err = api.parseRateResponse(res, currency)
	if err != nil {
		return currentRate, err
	}

	return currentRate, nil
}

func (api *HttpRateProvider) executeRateRequest(currency *models.Currency) (*http.Response, error) {
	req, err := api.executor.GenerateHttpRequest(currency)
	if err != nil {
		logger.Log.Error(err)
		return nil, errors.Wrap(err, "executeRateRequest")
	}
	res, err := api.client.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.Wrap(ErrFailedToGetRate, "executeRateRequest")
	}
	return res, nil
}

func (api *HttpRateProvider) parseRateResponse(response *http.Response, currency *models.Currency) (models.Rate, error) {
	currentRate, err := api.executor.ExtractRate(response, currency)
	if err != nil {
		logger.Log.Errorf("failed to extract rate, provider: %v", api.executor.Name())
		return currentRate, errors.Wrap(err, "parseRateResponse")
	}

	return currentRate, nil
}

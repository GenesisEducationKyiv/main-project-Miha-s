package rateProviders

import (
	"btc-test-task/internal/helpers/config"
	errors2 "btc-test-task/internal/helpers/customErrors"
	"btc-test-task/internal/helpers/models"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/pkg/errors"
)

type CoinAPIExecutor struct {
	endpoint string
	apiKey   string
}

func NewCoinAPIExecutor(conf *config.Config) *CoinAPIExecutor {
	return &CoinAPIExecutor{
		endpoint: conf.CoinAPIUrl,
		apiKey:   conf.CoinAPIKey,
	}
}

func (api *CoinAPIExecutor) ExtractRate(jsonValue []byte, _ *models.Currency) (models.Rate, error) {
	rate := models.Rate{}

	var dat map[string]interface{}
	if err := json.Unmarshal(jsonValue, &dat); err != nil {
		return rate, errors.Wrap(errors2.ErrFailedToGetRate, "extractRate")
	}
	var ok bool
	rate.Value, ok = dat["rate"].(float64)
	if !ok {
		return rate, errors.Wrap(errors2.ErrFailedToGetRate, "extractRate")
	}
	return rate, nil
}

func (api *CoinAPIExecutor) generateEndpoint(currency *models.Currency) (string, error) {
	return fmt.Sprintf(api.endpoint, currency.CurrencyFrom, currency.CurrencyTo), nil
}

func (api *CoinAPIExecutor) GenerateHttpRequest(currency *models.Currency) (*http.Request, error) {
	endpoint, err := api.generateEndpoint(currency)
	if err != nil {
		return nil, errors.Wrap(err, "generateCurrencyRequest")
	}
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(errors2.ErrFailedToGetRate, "generateCurrencyRequest")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CoinAPI-Key", api.apiKey)

	return req, nil
}

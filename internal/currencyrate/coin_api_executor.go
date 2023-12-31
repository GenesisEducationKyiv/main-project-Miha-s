package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/models"
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

func (api *CoinAPIExecutor) GenerateHttpRequest(currency *models.Currency) (*http.Request, error) {
	endpoint := api.generateEndpoint(currency)
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(ErrFailedToGetRate, "GenerateHttpRequest")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CoinAPI-Key", api.apiKey)

	return req, nil
}

func (api *CoinAPIExecutor) ExtractRate(resp *http.Response, _ *models.Currency) (models.Rate, error) {
	rate := struct {
		Rate float64 `json:"rate"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&rate)
	if err != nil {
		return models.Rate{}, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}

	return models.Rate{
		Value: rate.Rate,
	}, nil
}

func (api *CoinAPIExecutor) generateEndpoint(currency *models.Currency) string {
	return fmt.Sprintf(api.endpoint, currency.From, currency.To)
}

func (api *CoinAPIExecutor) Name() string {
	return "CoinAPI"
}

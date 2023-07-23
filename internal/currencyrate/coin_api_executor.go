package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	models2 "btc-test-task/internal/common/models"
	"bytes"
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

func (api *CoinAPIExecutor) GenerateHttpRequest(currency *models2.Currency) (*http.Request, error) {
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

func (api *CoinAPIExecutor) ExtractRate(jsonValue []byte, _ *models2.Currency) (models2.Rate, error) {
	rate := struct {
		Rate float64 `json:"rate"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(jsonValue))
	err := decoder.Decode(&rate)
	if err != nil {
		return models2.Rate{}, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}

	return models2.Rate{
		Value: rate.Rate,
	}, nil
}

func (api *CoinAPIExecutor) generateEndpoint(currency *models2.Currency) string {
	return fmt.Sprintf(api.endpoint, currency.From, currency.To)
}

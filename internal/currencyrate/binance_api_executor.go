package currencyrate

import (
	"btc-test-task/internal/configuration/config"
	models "btc-test-task/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type BinanceAPIExecutor struct {
	endpoint string
}

func NewBinanceAPIExecutor(conf *config.Config) *BinanceAPIExecutor {
	return &BinanceAPIExecutor{
		endpoint: conf.BinanceAPIUrl,
	}
}

func (api *BinanceAPIExecutor) GenerateHttpRequest(currency *models.Currency) (*http.Request, error) {
	endpoint := api.generateEndpoint(currency)
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)

	if err != nil {
		return nil, errors.Wrap(ErrFailedToGetRate, "GenerateHttpRequest")
	}

	return req, nil
}

func (api *BinanceAPIExecutor) ExtractRate(jsonValue []byte, _ *models.Currency) (models.Rate, error) {
	price := struct {
		Price string `json:"price"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(jsonValue))
	err := decoder.Decode(&price)
	if err != nil {
		return models.Rate{}, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}
	rate, err := strconv.ParseFloat(price.Price, 64)
	if err != nil {
		return models.Rate{}, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}

	return models.Rate{
		Value: rate,
	}, nil
}

func (api *BinanceAPIExecutor) generateEndpoint(currency *models.Currency) string {
	return fmt.Sprintf(api.endpoint, currency.CurrencyFrom, currency.CurrencyTo)
}

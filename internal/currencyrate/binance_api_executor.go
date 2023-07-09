package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	models2 "btc-test-task/internal/common/models"
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

func (api *BinanceAPIExecutor) GenerateHttpRequest(currency *models2.Currency) (*http.Request, error) {
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

func (api *BinanceAPIExecutor) ExtractRate(jsonValue []byte, _ *models2.Currency) (models2.Rate, error) {
	price := struct {
		Price string `json:"price"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(jsonValue))
	err := decoder.Decode(&price)
	if err != nil {
		return models2.Rate{}, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}
	rate, err := strconv.ParseFloat(price.Price, 64)
	if err != nil {
		return models2.Rate{}, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}

	return models2.Rate{
		Value: rate,
	}, nil
}

func (api *BinanceAPIExecutor) generateEndpoint(currency *models2.Currency) string {
	return fmt.Sprintf(api.endpoint, currency.From, currency.To)
}

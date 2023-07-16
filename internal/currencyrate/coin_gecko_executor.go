package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/models"
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

func (api *CoinGeckoExecutor) GenerateHttpRequest(currency *models.Currency) (*http.Request, error) {
	endpoint := api.generateEndpoint(currency)
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "GenerateHttpRequest")
	}
	return req, err
}

func (api *CoinGeckoExecutor) ExtractRate(resp *http.Response, currency *models.Currency) (models.Rate, error) {
	rate := models.Rate{}

	var responseData CoinGeckoResponse

	err := json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return rate, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}

	rate.Value = responseData[currenciesCoinGecko[currency.From]][currenciesCoinGecko[currency.To]]
	if rate.Value == 0 {
		return rate, errors.Wrap(ErrFailedToGetRate, "ExtractRate")
	}
	return rate, nil
}

func (api *CoinGeckoExecutor) generateEndpoint(currency *models.Currency) string {
	return fmt.Sprintf(api.url, currency.From, currency.To)
}

func (api *CoinGeckoExecutor) Name() string {
	return "CoinGecko"
}

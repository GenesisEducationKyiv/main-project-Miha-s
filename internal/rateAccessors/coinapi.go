package rateAccessors

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
	"encoding/json"

	"io"
	"net/http"

	"github.com/pkg/errors"
)

type CoinAPI struct {
	endpoint string
	apiKey   string
}

func NewCoinAPI(conf *config.Config) (*CoinAPI, error) {
	newCoinAPI := new(CoinAPI)
	err := newCoinAPI.init(conf)
	if err != nil {
		return nil, err
	}
	return newCoinAPI, nil
}

func (api *CoinAPI) init(conf *config.Config) error {
	api.endpoint = conf.CoinAPIUrl + conf.CurrencyFrom + "/" + conf.CurrencyTo
	api.apiKey = conf.CoinAPIKey
	return nil
}

func extractRate(resp *http.Response) (float64, error) {
	jsonValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(ErrFailedToGetRate, "extractRate")
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonValue, &dat); err != nil {
		return 0, errors.Wrap(ErrFailedToGetRate, "extractRate")
	}
	rate, ok := dat["rate"].(float64)
	if !ok {
		return 0, errors.Wrap(ErrFailedToGetRate, "extractRate")
	}
	return rate, nil
}

func (api *CoinAPI) generateEndpoint(currencyFrom string, currencyTo string) (string, error) {
	return api.endpoint + currencyFrom + "/" + currencyTo, nil
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

func (api *CoinAPI) GetCurrentRate(currencyFrom string, currencyTo string) (float64, error) {
	value := 0.0
	req, err := api.generateCurrencyRequest(currencyFrom, currencyTo)
	if err != nil {
		logger.Log.Error(err)
		return value, errors.Wrap(err, "GetCurrentRate")
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return value, errors.Wrap(ErrFailedToGetRate, "GetCurrentRate")
	}

	value, err = extractRate(res)
	if err != nil {
		return value, errors.Wrap(err, "GetCurrentRate")
	}

	logger.Log.Infof("The rate %v", value)
	return value, nil
}

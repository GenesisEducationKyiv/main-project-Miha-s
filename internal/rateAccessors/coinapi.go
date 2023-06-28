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

func extractRate(jsonValue []byte) (float64, error) {
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

func (api *CoinAPI) GetCurrentRate() (float64, error) {
	value := 0.0
	req, err := http.NewRequest(
		http.MethodGet,
		api.endpoint,
		nil,
	)

	if err != nil {
		return value, errors.Wrap(ErrFailedToGetRate, "GetCurrentRate")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CoinAPI-Key", api.apiKey)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return value, errors.Wrap(ErrFailedToGetRate, "GetCurrentRate")
	}
	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return value, errors.Wrap(ErrFailedToGetRate, "GetCurrentRate")
	}
	value, err = extractRate(responseBytes)
	if err != nil {
		return value, errors.Wrap(err, "GetCurrentRate")
	}

	logger.Log.Infof("The rate %v", value)
	return value, nil
}

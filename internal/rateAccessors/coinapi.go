package rateAccessors

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CoinApI struct {
	endpoint string
	apiKey   string
}

func (api *CoinApI) Init(conf *config.Config) error {
	api.endpoint = conf.CoinAPIUrl + conf.CurrencyFrom + "/" + conf.CurrencyTo
	api.apiKey = conf.CoinAPIKey
	return nil
}

func extractRate(jsonValue []byte) (float64, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonValue, &dat); err != nil {
		return 0, err
	}
	rate, ok := dat["rate"].(float64)
	if !ok {
		return 0, errors.New("cannot extract float rate value")
	}
	return rate, nil
}

func (api *CoinApI) GetCurrentRate() (float64, error) {
	value := 0.0
	req, err := http.NewRequest(
		http.MethodGet,
		api.endpoint,
		nil,
	)

	if err != nil {
		return value, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CoinAPI-Key", api.apiKey)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return value, err
	}
	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return value, err
	}
	value, err = extractRate(responseBytes)
	if err != nil {
		return value, err
	}

	logger.LogInfo(fmt.Sprintf("The rate %v", value))
	return value, err
}

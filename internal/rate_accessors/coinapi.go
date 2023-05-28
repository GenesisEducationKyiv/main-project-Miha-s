package rate_accessors

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CoinApI struct {
	endpoint string
	api_key  string
}

func (api *CoinApI) Init(conf *config.Config) error {
	api.endpoint = conf.CoinAPIUrl + conf.CurrencyFrom + "/" + conf.CurrencyTo
	api.api_key = conf.CoinAPIKey
	return nil
}

func extract_rate(json_value []byte) (float64, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(json_value, &dat); err != nil {
		return 0, err
	}
	return dat["rate"].(float64), nil
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
	req.Header.Add("X-CoinAPI-Key", api.api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return value, err
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	value, err = extract_rate(responseBytes)
	logger.LogInfo(fmt.Sprintf("The rate %v", value))
	return value, err
}

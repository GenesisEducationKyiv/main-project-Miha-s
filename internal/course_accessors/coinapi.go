package course_accessors

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
}

func (api *CoinApI) Init(conf *config.Config) {
	api.endpoint = conf.CoinAPIUrl
}

func extract_course(json_value []byte) (float64, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(json_value, &dat); err != nil {
		return 0, err
	}
	return dat["rate"].(float64), nil
}

func (api *CoinApI) GetBTCToUAHCourse() (float64, error) {
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
	req.Header.Add("X-CoinAPI-Key", "C72B3E65-2D6F-4AFC-AA94-AD84537D1825")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return value, err
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	value, err = extract_course(responseBytes)
	logger.LogInfo(fmt.Sprintf("The course %v", value))
	return value, err
}

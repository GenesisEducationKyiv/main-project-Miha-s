package course_accessors

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type coinAPI struct {
	endpoint string
}

var coin_api coinAPI

func Init(conf *config.Config) {

}

func extract_course(json_value []byte) (float64, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(json_value, &dat); err != nil {
		return 0, err
	}
	return dat["rate"].(float64), nil
}

func GetBTCToUAHCourse() (float64, error) {
	value := 0.0
	req, err := http.NewRequest(
		http.MethodGet,
		"https://rest.coinapi.io/v1/exchangerate/BTC/UAH",
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

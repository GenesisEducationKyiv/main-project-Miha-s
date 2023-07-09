package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestHttpRateProvider(t *testing.T) {
	err := logger.Init(&config.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	t.Run("SuccessfulRequestResponse", testSuccessfulRequestResponse)
	t.Run("BadRequest", testBadRequest)
	t.Run("BadRate", testBadRate)
}

func testSuccessfulRequestResponse(t *testing.T) {
	rate := models.Rate{Value: 123}
	provider := NewHttpRateProvider(&rateExecutorStub{
		Rate: rate,
	}, &httpClientStub{})

	res, err := provider.GetCurrentRate(&models.Currency{})
	if err != nil {
		t.Errorf("Failed to execute correct request: %v", err)
	}

	if res != rate {
		t.Errorf("Failed to get correct rate, got %v, expected %v", res, rate)
	}
}

func testBadRequest(t *testing.T) {
	provider := NewHttpRateProvider(&rateExecutorStub{
		RequestErr: errors.New("Some error"),
	}, &httpClientStub{})

	_, err := provider.GetCurrentRate(&models.Currency{})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func testBadRate(t *testing.T) {
	provider := NewHttpRateProvider(&rateExecutorStub{
		RateErr: errors.New("Rate error"),
	}, &httpClientStub{})

	_, err := provider.GetCurrentRate(&models.Currency{})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

type httpClientStub struct{}

func (client *httpClientStub) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       io.NopCloser(strings.NewReader("asdf")),
		StatusCode: http.StatusOK,
	}, nil
}

type rateExecutorStub struct {
	Request    http.Request
	RequestErr error
	Rate       models.Rate
	RateErr    error
}

func (executor *rateExecutorStub) GenerateHttpRequest(_ *models.Currency) (*http.Request, error) {
	return &executor.Request, executor.RequestErr
}
func (executor *rateExecutorStub) ExtractRate(_ []byte, _ *models.Currency) (models.Rate, error) {
	return executor.Rate, executor.RateErr
}

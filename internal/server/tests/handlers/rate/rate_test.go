package rateTest

import (
	"btc-test-task/internal/configuration/config"
	"btc-test-task/internal/configuration/logger"
	"btc-test-task/internal/currencyrate"
	rateAccessorsTest "btc-test-task/internal/currencyrate/tests"
	templatesTest "btc-test-task/internal/email/templates/tests"
	"btc-test-task/internal/lifecycle"
	"btc-test-task/internal/models"
	"btc-test-task/internal/server/handlers"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var conf config.Config

func globalSetup() error {
	err := logger.Init(&conf)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	err := globalSetup()
	if err != nil {
		panic(fmt.Sprintf("failed to setup test %v", err))
	}
	code := m.Run()

	os.Exit(code)
}

func createRateRequest() *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/api/rate", bytes.NewReader([]byte("")))
	return request
}

func createRateHandler(rateAccessor handlers.RateProvider) (http.HandlerFunc, error) {
	servicesStubs := new(lifecycle.Services)
	servicesStubs.Templates = &templatesTest.TemplatesImplStub{}
	servicesStubs.RateProvider = rateAccessor
	factory := handlers.NewHandlersFactoryImpl(&conf, servicesStubs)
	rateHandler := factory.CreateRate()
	return rateHandler, nil
}

func TestGetValidRate(t *testing.T) {
	rateValue := models.Rate{Value: 3849.4}
	rateHandler, err := createRateHandler(&rateAccessorsTest.RateProviderStub{
		Rate: rateValue,
	})
	if err != nil {
		t.Fatalf("failed to create rate handler %v", err)
	}
	rateRequest := createRateRequest()
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusOK

	rateHandler(writer, rateRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

func TestGetInvalidRate(t *testing.T) {
	rateHandler, err := createRateHandler(&rateAccessorsTest.RateProviderStub{
		RateError: currencyrate.ErrFailedToGetRate,
	})
	if err != nil {
		t.Fatalf("failed to create rate handler %v", err)
	}
	rateRequest := createRateRequest()
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusNotFound

	rateHandler(writer, rateRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

package handlers

import (
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"btc-test-task/internal/currencyrate"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateHandler(t *testing.T) {
	err := logger.Init(&conf)
	if err != nil {
		t.Fatalf("Failed to initilize logger: %v", err)
	}
	t.Run("GetValidRate", testGetValidRate)
	t.Run("GetInvalidRate", testGetInvalidRate)
}

func createRateRequest() *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/api/rate", bytes.NewReader([]byte("")))
	return request
}

func createRateHandler(rateAccessor RateProvider) (http.HandlerFunc, error) {
	servicesStubs := new(ServicesStub)
	servicesStubs.Templates = &templatesImplStub{}
	servicesStubs.RateProvider = rateAccessor
	factory := NewHandlersFactoryImpl(&conf, servicesStubs)
	rateHandler := factory.CreateRate()
	return rateHandler, nil
}

func testGetValidRate(t *testing.T) {
	rateValue := models.Rate{Value: 3849.4}
	rateHandler, err := createRateHandler(&rateProviderStub{
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

func testGetInvalidRate(t *testing.T) {
	rateHandler, err := createRateHandler(&rateProviderStub{
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

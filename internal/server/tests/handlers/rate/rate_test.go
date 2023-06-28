package rateTest

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
	templatesTest "btc-test-task/internal/helpers/templates/tests"
	"btc-test-task/internal/helpers/types"
	"btc-test-task/internal/rateAccessors"
	rateAccessorsTest "btc-test-task/internal/rateAccessors/tests"
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
	err := conf.LoadFromENV("../../../../../.env.test")
	if err != nil {
		return err
	}
	err = logger.Init(&conf)
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

func createRateHandler(rateAccessor rateAccessors.RateAccessor) (http.HandlerFunc, error) {
	servicesStubs := new(types.Services)
	servicesStubs.Templates = &templatesTest.TemplatesImplStub{}
	servicesStubs.RateAccessor = rateAccessor
	factory, err := handlers.NewHandlersFactoryImpl(&conf, servicesStubs)
	if err != nil {
		return nil, err
	}
	rateHandler := factory.CreateRate()
	return rateHandler, nil
}

func TestGetValidRate(t *testing.T) {
	rateValue := 3849.4
	rateHandler, err := createRateHandler(&rateAccessorsTest.RateAccessorStub{
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
	rateHandler, err := createRateHandler(&rateAccessorsTest.RateAccessorStub{
		RateError: rateAccessors.ErrFailedToGetRate,
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

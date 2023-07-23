package handlers

import (
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubscribe(t *testing.T) {
	err := logger.Init(&conf)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	t.Run("SubscribeValidRequest", testSubscribeValidRequest)
	t.Run("SubscribeInvalidRequest", testSubscribeInvalidRequest)
	t.Run("SubscribeExistingEmail", testSubscribeExistingEmail)
}

func createSubscribeBody(email models.Email) io.Reader {
	return bytes.NewReader([]byte(fmt.Sprintf("email=%v", email.Value)))
}

func createSubscribeRequest(body io.Reader) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/api/subscribe", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func createSubscribeHandler() (http.HandlerFunc, error) {
	servicesStubs := new(ServicesStub)
	servicesStubs.EmailsRepository = &emailsStorageStub{}
	factory := NewHandlersFactoryImpl(&conf, servicesStubs)
	subscribeHandler := factory.CreateSubscribe()
	return subscribeHandler, nil
}

func testSubscribeValidRequest(t *testing.T) {
	subscribeHandler, err := createSubscribeHandler()
	if err != nil {
		t.Fatalf("failed to create subscribe handler %v", err)
	}
	body := createSubscribeBody(GoodEmail)
	goodSubscribeRequest := createSubscribeRequest(body)
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusOK

	subscribeHandler(writer, goodSubscribeRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

func testSubscribeInvalidRequest(t *testing.T) {
	subscribeHandler, err := createSubscribeHandler()
	if err != nil {
		t.Fatalf("failed to create subscribe handler %v", err)
	}
	body := createSubscribeBody(BadEmail)
	badSubscribeRequest := createSubscribeRequest(body)
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusBadRequest

	subscribeHandler(writer, badSubscribeRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

func testSubscribeExistingEmail(t *testing.T) {
	subscribeHandler, err := createSubscribeHandler()
	if err != nil {
		t.Fatalf("failed to create subscribe handler %v", err)
	}
	body := createSubscribeBody(ExistingEmail)
	badSubscribeRequest := createSubscribeRequest(body)
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusConflict

	subscribeHandler(writer, badSubscribeRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

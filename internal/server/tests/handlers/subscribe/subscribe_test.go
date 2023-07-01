package subscribeTest

import (
	emailsStorageTest "btc-test-task/internal/emailsRepository/tests"
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/models"
	"btc-test-task/internal/helpers/types"
	"btc-test-task/internal/server/handlers"
	"bytes"
	"fmt"
	"io"
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

func createSubscribeBody(email models.Email) io.Reader {
	return bytes.NewReader([]byte(fmt.Sprintf("email=%v", email.Value)))
}

func createSubscribeRequest(body io.Reader) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/api/subscribe", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func createSubscribeHandler() (http.HandlerFunc, error) {
	servicesStubs := new(types.Services)
	servicesStubs.EmailsRepository = &emailsStorageTest.EmailsStorageStub{}
	factory, err := handlers.NewHandlersFactoryImpl(&conf, servicesStubs)
	if err != nil {
		return nil, err
	}
	subscribeHandler := factory.CreateSubscribe()
	return subscribeHandler, nil
}

func TestSubscribeValidRequest(t *testing.T) {
	subscribeHandler, err := createSubscribeHandler()
	if err != nil {
		t.Fatalf("failed to create subscribe handler %v", err)
	}
	body := createSubscribeBody(emailsStorageTest.GoodEmail)
	goodSubscribeRequest := createSubscribeRequest(body)
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusOK

	subscribeHandler(writer, goodSubscribeRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

func TestSubscribeInvalidRequest(t *testing.T) {
	subscribeHandler, err := createSubscribeHandler()
	if err != nil {
		t.Fatalf("failed to create subscribe handler %v", err)
	}
	body := createSubscribeBody(emailsStorageTest.BadEmail)
	badSubscribeRequest := createSubscribeRequest(body)
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusBadRequest

	subscribeHandler(writer, badSubscribeRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

func TestSubscribeExistingEmail(t *testing.T) {
	subscribeHandler, err := createSubscribeHandler()
	if err != nil {
		t.Fatalf("failed to create subscribe handler %v", err)
	}
	body := createSubscribeBody(emailsStorageTest.ExistingEmail)
	badSubscribeRequest := createSubscribeRequest(body)
	writer := httptest.NewRecorder()
	expectedStatus := http.StatusConflict

	subscribeHandler(writer, badSubscribeRequest)

	if writer.Code != expectedStatus {
		t.Fatalf("expected status %v, got %v", expectedStatus, writer.Code)
	}
}

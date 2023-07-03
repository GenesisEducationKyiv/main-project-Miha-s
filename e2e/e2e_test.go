package e2e

import (
	"btc-test-task/internal/configuration/config"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"
)

var conf config.Config

func globalSetup() error {
	err := conf.LoadFromENV("../.env")
	return err
}

func TestMain(m *testing.M) {
	err := globalSetup()
	if err != nil {
		os.Exit(2)
	}
	code := m.Run()
	os.Exit(code)
}

func TestSubscribe(t *testing.T) {
	expectedResponseStatus := http.StatusOK
	emailToSubscribe := "someemail@gmail.com"
	resp, err := http.Post(fmt.Sprintf("http://localhost:%v/api/subscribe", conf.Port),
		"application/x-www-form-urlencoded", bytes.NewReader([]byte("email="+emailToSubscribe)))
	defer resp.Body.Close()

	if err != nil {
		t.Fatalf("failed to send request %v", err)
	}

	if resp.StatusCode != expectedResponseStatus {
		t.Fatalf("expected status code %v, got %v", expectedResponseStatus, resp.StatusCode)
	}
}

func TestGetRate(t *testing.T) {
	expectedResponseStatus := http.StatusOK
	resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/rate", conf.Port))
	defer resp.Body.Close()

	if err != nil {
		t.Fatalf("failed to send request %v", err)
	}

	if resp.StatusCode != expectedResponseStatus {
		t.Fatalf("expected status code %v, got %v", expectedResponseStatus, resp.StatusCode)
	}
}

package handlers

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/types"

	"net/http"
)

type HandlersFactory interface {
	Init(conf *config.Config, services *types.Services) error
	CreateRate() http.HandlerFunc
	CreateSubscribe() http.HandlerFunc
	CreateSendEmails() http.HandlerFunc
}

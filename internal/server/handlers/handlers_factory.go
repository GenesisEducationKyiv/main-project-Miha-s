package handlers

import (
	"net/http"
)

type HandlersFactory interface {
	CreateRate() http.HandlerFunc
	CreateSubscribe() http.HandlerFunc
	CreateSendEmails() http.HandlerFunc
}

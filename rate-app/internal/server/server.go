package server

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlersFactory interface {
	CreateRate() http.HandlerFunc
	CreateSubscribe() http.HandlerFunc
	CreateSendEmails() http.HandlerFunc
}

type Server struct {
	httpServer *http.Server
}

func NewServer(conf *config.Config, handlersFactory HandlersFactory, logger MiddlewareFunc) *Server {
	newServer := new(Server)
	newServer.init(conf, handlersFactory, logger)

	return newServer
}

func (serv *Server) init(conf *config.Config, handlersFactory HandlersFactory, logger MiddlewareFunc) {
	router := chi.NewRouter()
	router.Use(logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Get("/rate", handlersFactory.CreateRate())
		r.Post("/subscribe", handlersFactory.CreateSubscribe())
		r.Post("/sendEmails", handlersFactory.CreateSendEmails())
	})

	handler := router

	serv.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Port),
		Handler: handler,
	}
}

func (serv *Server) Run() error {
	err := serv.httpServer.ListenAndServe()
	return err
}

func (serv *Server) Shutdown() {
	err := serv.httpServer.Shutdown(context.Background())
	if err != nil {
		logger.Log.Errorf("Failed to shutdown server %v", err)
	}
}

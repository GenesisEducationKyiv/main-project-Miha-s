package server

import (
	"btc-test-task/internal/configuration/config"
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
	router chi.Router
	port   uint
}

func NewServer(conf *config.Config, handlersFactory HandlersFactory) *Server {
	newServer := new(Server)
	newServer.init(conf, handlersFactory)

	return newServer
}

func (serv *Server) init(conf *config.Config, handlersFactory HandlersFactory) {
	serv.port = conf.Port

	serv.router = chi.NewRouter()
	serv.router.Use(middleware.RequestID)
	serv.router.Use(middleware.RealIP)
	serv.router.Use(middleware.Logger)
	serv.router.Use(middleware.Recoverer)

	serv.router.Use(middleware.Timeout(60 * time.Second))

	serv.router.Route("/api", func(r chi.Router) {
		r.Get("/rate", handlersFactory.CreateRate())
		r.Post("/subscribe", handlersFactory.CreateSubscribe())
		r.Post("/sendEmails", handlersFactory.CreateSendEmails())
	})
}

func (serv *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%v", serv.port), serv.router)
	return err
}

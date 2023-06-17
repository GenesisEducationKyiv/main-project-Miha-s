package server

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/handlers"
	"fmt"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router chi.Router
	port   uint
}

func (serv *Server) Init(conf *config.Config, handlers_factory handlers.HandlersFactory) error {
	serv.port = conf.Port

	serv.router = chi.NewRouter()
	serv.router.Use(middleware.RequestID)
	serv.router.Use(middleware.RealIP)
	serv.router.Use(middleware.Logger)
	serv.router.Use(middleware.Recoverer)

	serv.router.Use(middleware.Timeout(60 * time.Second))

	serv.router.Route("/api", func(r chi.Router) {
		r.Get("/rate", handlers_factory.CreateRate())
		r.Post("/subscribe", handlers_factory.CreateSubscribe())
		r.Post("/sendEmails", handlers_factory.CreateSendEmails())
	})

	return nil
}

func (serv *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%v", serv.port), serv.router)
	return err
}

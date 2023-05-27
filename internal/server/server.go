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

func handle_root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Some other stuff"))
	return
}

type Server struct {
	router chi.Router
	port   uint
}

func (serv *Server) Init(conf *config.Config) {
	serv.router = chi.NewRouter()
	serv.router.Use(middleware.RequestID)
	serv.router.Use(middleware.RealIP)
	serv.router.Use(middleware.Logger)
	serv.router.Use(middleware.Recoverer)

	serv.router.Use(middleware.Timeout(60 * time.Second))

	serv.router.Get("/rate", handlers.Rate)
	serv.router.Post("/subscribe", handlers.Subscribe)
	serv.router.Post("/sendEmails", handlers.SendEmails)
}

func (serv *Server) Run() {

	err := http.ListenAndServe(":2777", serv.router)
	fmt.Println(err)
}

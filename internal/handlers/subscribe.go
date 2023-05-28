package handlers

import (
	"btc-test-task/internal/logger"
	"net/http"
)

func (factory *HandlersFactoryImpl) CreateSubscribe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		err := factory.services.EmailStorage.AddEmail(email)
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

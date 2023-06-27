package handlers

import (
	"btc-test-task/internal/emailsStorage"
	"btc-test-task/internal/helpers/logger"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateSubscribe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		valid := factory.services.EmailStorage.ValidateEmail(email)
		if !valid {
			logger.Log.Info("Incorrect email")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := factory.services.EmailStorage.AddEmail(email)
		if errors.Is(err, emailsStorage.ErrEmailAlreadyExists) {
			logger.Log.Info(err)
			w.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

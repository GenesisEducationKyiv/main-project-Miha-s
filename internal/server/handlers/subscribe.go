package handlers

import (
	errors2 "btc-test-task/internal/helpers/errors"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/models"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateSubscribe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emailStr := r.FormValue("email")
		email := models.Email{Value: emailStr}

		err := factory.services.GetEmailsRepositoryService().AddEmail(&email)

		if errors.Is(err, errors2.ErrInvalidEmailAddress) {
			logger.Log.Info(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, errors2.ErrEmailAlreadyExists) {
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

package handlers

import (
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"btc-test-task/internal/repository"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateSubscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailStr := r.FormValue("email")
		email := models.Email{Value: emailStr}

		err := factory.services.GetEmailsRepositoryService().AddEmail(&email)

		if errors.Is(err, repository.ErrInvalidEmailAddress) {
			logger.Log.Info(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, repository.ErrEmailAlreadyExists) {
			logger.Log.Info(err)
			w.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

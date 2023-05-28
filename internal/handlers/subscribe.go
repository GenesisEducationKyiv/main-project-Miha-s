package handlers

import (
	"btc-test-task/internal/logger"
	"net/http"

	"github.com/go-chi/render"
)

func (factory *HandlersFactoryImpl) CreateSubscribe() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		logger.LogInfo(email)
		err := factory.services.EmailStorage.AddEmail(email)
		if err != nil {
			logger.LogError(err)
			render.Status(r, http.StatusConflict)
			return
		}

		render.Status(r, http.StatusOK)
	})
}

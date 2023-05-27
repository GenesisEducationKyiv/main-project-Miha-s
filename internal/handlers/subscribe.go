package handlers

import (
	"btc-test-task/internal/emails_storage"
	"btc-test-task/internal/logger"
	"net/http"

	"github.com/go-chi/render"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	logger.LogInfo(email)
	err := emails_storage.AddEmail(email)
	if err != nil {
		logger.LogError(err)
		render.Status(r, http.StatusConflict)
		return
	}

	render.Status(r, http.StatusOK)
}

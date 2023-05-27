package handlers

import (
	"btc-test-task/internal/course_accessors"
	"btc-test-task/internal/email_sender"
	"btc-test-task/internal/emails_storage"
	"fmt"

	"btc-test-task/internal/logger"
	"net/http"

	"github.com/go-chi/render"
)

func SendEmails(w http.ResponseWriter, r *http.Request) {
	emails := emails_storage.GetAllEmails()
	course, err := course_accessors.GetBTCToUAHCourse()
	if err != nil {
		logger.LogError(err)
		render.Status(r, http.StatusInternalServerError)
		return
	}

	for email := range *emails {
		email_sender.SendEmail(email, fmt.Sprint(course))
	}

	render.Status(r, http.StatusOK)
}

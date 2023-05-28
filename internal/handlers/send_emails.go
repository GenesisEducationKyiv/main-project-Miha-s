package handlers

import (
	"btc-test-task/internal/logger"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (factory *HandlersFactoryImpl) CreateSendEmails() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emails := factory.services.EmailStorage.GetAllEmails()
		course, err := factory.services.CourseAccessor.GetBTCToUAHCourse()
		if err != nil {
			logger.LogError(err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		for email := range *emails {
			factory.services.EmailSender.SendEmail(email, fmt.Sprint(course))
		}

		render.Status(r, http.StatusOK)
	})
}

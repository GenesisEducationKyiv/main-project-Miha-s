package handlers

import (
	"btc-test-task/internal/logger"
	"net/http"
)

func (factory *HandlersFactoryImpl) CreateSendEmails() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emails := factory.services.EmailStorage.GetAllEmails()
		rate, err := factory.services.RateAccessor.GetCurrentRate()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go factory.services.EmailSender.BroadcastEmails(emails, factory.services.Templates.CurrencyRate(rate))

		w.WriteHeader(http.StatusOK)
	})
}

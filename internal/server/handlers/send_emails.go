package handlers

import (
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/rateProviders"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateSendEmails() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emails := factory.services.EmailStorage.GetAllEmails()
		rate, err := factory.services.RateProvider.GetCurrentRate(factory.currencyFrom, factory.currencyTo)
		if errors.Is(err, rateProviders.ErrFailedToGetRate) {
			logger.Log.Warn(err)
			w.WriteHeader(http.StatusFailedDependency)
			return
		} else if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go factory.services.EmailSender.BroadcastEmails(emails, factory.services.Templates.CurrencyRate(rate))

		w.WriteHeader(http.StatusOK)
	})
}

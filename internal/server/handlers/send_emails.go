package handlers

import (
	errors2 "btc-test-task/internal/helpers/errors"
	"btc-test-task/internal/helpers/logger"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateSendEmails() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emails := factory.services.GetEmailsRepositoryService().GetAllEmails()
		rate, err := factory.services.GetRateProviderService().GetCurrentRate(factory.currencyFrom, factory.currencyTo)
		if errors.Is(err, errors2.ErrFailedToGetRate) {
			logger.Log.Warn(err)
			w.WriteHeader(http.StatusFailedDependency)
			return
		} else if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go factory.services.GetEmailSenderService().BroadcastEmails(
			emails, factory.services.GetTemplatesService().CurrencyRate(rate))

		w.WriteHeader(http.StatusOK)
	})
}

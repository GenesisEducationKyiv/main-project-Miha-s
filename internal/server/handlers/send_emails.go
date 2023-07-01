package handlers

import (
	errors2 "btc-test-task/internal/helpers/customErrors"
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/helpers/models"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateSendEmails() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emails := factory.services.GetEmailsRepositoryService().GetAllEmails()
		rate, err := factory.services.GetRateProviderService().GetCurrentRate(
			&models.Currency{CurrencyFrom: factory.currencyFrom, CurrencyTo: factory.currencyTo})
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

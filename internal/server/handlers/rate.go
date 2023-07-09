package handlers

import (
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"btc-test-task/internal/currencyrate"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateRate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rate, err := factory.services.GetRateProviderService().GetCurrentRate(
			&models.Currency{From: factory.currencyFrom, To: factory.currencyTo})
		if errors.Is(err, currencyrate.ErrFailedToGetRate) {
			logger.Log.Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte(factory.services.GetTemplatesService().CurrencyRate(rate)))
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

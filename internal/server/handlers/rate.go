package handlers

import (
	errors2 "btc-test-task/internal/helpers/errors"
	"btc-test-task/internal/helpers/logger"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateRate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rate, err := factory.services.GetRateProviderService().GetCurrentRate(factory.currencyFrom, factory.currencyTo)
		if errors.Is(err, errors2.ErrFailedToGetRate) {
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
	})
}

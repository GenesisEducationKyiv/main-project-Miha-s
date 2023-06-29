package handlers

import (
	"btc-test-task/internal/helpers/logger"
	"btc-test-task/internal/rateAccessors"
	"net/http"

	"github.com/pkg/errors"
)

func (factory *HandlersFactoryImpl) CreateRate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("some info")
		rate, err := factory.services.RateAccessor.GetCurrentRate(factory.currencyFrom, factory.currencyTo)
		if errors.Is(err, rateAccessors.ErrFailedToGetRate) {
			logger.Log.Warn(err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte(factory.services.Templates.CurrencyRate(rate)))
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

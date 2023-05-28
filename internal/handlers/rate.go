package handlers

import (
	"btc-test-task/internal/logger"
	"net/http"
)

func (factory *HandlersFactoryImpl) CreateRate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rate, err := factory.services.RateAccessor.GetCurrentRate()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(factory.services.Templates.CurrencyRate(rate)))
		w.WriteHeader(http.StatusOK)
	})
}

package handlers

import (
	"btc-test-task/internal/logger"
	"fmt"
	"net/http"
)

func (factory *HandlersFactoryImpl) CreateRate() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		course, err := factory.services.CourseAccessor.GetBTCToUAHCourse()
		if err != nil {
			logger.LogError(err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprintf("%v", course)))
	})
}

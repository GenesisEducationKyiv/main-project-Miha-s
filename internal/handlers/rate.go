package handlers

import (
	"btc-test-task/internal/course_accessors"
	"btc-test-task/internal/logger"
	"fmt"
	"net/http"
)

func Rate(w http.ResponseWriter, r *http.Request) {
	course, err := course_accessors.GetBTCToUAHCourse()
	if err != nil {
		logger.LogError(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write([]byte(fmt.Sprintf("%v", course)))
}

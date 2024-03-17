package services

import (
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/logout" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /logout")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected POST for URL: " + r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(status) +
				" | Details: " + string(err[:]))

		w.WriteHeader(status)
		w.Write(err)
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK) + " | Detail: token removed")
	r.Header.Del("Authorization")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("token removed"))
}

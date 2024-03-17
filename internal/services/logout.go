package services

import (
	"encoding/json"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"net/http"
)

// @Summary Logout
// @Description This endpoint for Logout
// @Produce json
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /logout [post]
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
		if status == http.StatusInternalServerError {
			logger.WarningLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
		} else {
			logger.InfoLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
		}

		jsonResp, err := json.Marshal(err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}
		w.WriteHeader(status)
		w.Write(jsonResp)
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK) + " | Detail: token removed")
	r.Header.Del("Authorization")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("token removed"))
}

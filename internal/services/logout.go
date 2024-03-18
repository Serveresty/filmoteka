package services

import (
	"encoding/json"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"net/http"
)

// @Summary Logout
// @Security ApiKeyAuth
// @Tags auth
// @Description This endpoint for Logout
// @Produce json
// @Success 200 {string} string "JSON с сообщением"
// @Failure 401 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/logout" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /logout")
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected POST for URL: " + r.URL.Path)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		if !IsTesting {
			if status == http.StatusInternalServerError {
				logger.WarningLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			} else {
				logger.InfoLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			}
		}

		jsonResp, err := json.Marshal(err)
		if err != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err.Error() +
						" | Details: " + "error while marshal array of errors")
			}

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

	var pull []models.ErrorInfo
	pull = append(pull, models.ErrorInfo{Type: "success", Message: "token removed"})
	jsonResp, err1 := json.Marshal(pull)
	if err1 != nil {
		if !IsTesting {
			logger.WarningLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusInternalServerError) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while marshal array of errors")
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal errs"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

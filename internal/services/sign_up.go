package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"net/http"
)

// @Summary SignUp Page
// @Description This endpoint for sign-up page
// @Produce json
// @Success 200
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /sign-up [get]
func SignUp(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/sign-up" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /sign-up")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected GET for URL: " + r.URL.Path)
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

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status))
	w.WriteHeader(status)
}

// @Summary Registration
// @Description This endpoint for registration
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /registration [post]
func Registration(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/registration" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /registration")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var user models.User

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

	err1 := json.NewDecoder(r.Body).Decode(&user)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error decoding json")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error decoding json"))
		return
	}

	if !jwt.ValidateAuthData(&user) {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Details: " + "wrong input data")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error wrong input data"))
		return
	}

	hashPass, err1 := jwt.HashPassword(user.Password)
	if err1 != nil {
		logger.WarningLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusInternalServerError) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while hashing password")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while hashing password"))
		return
	}

	user.Password = hashPass

	err1 = database.CreateNewUser(&user)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "user already created with same data")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error user already created with same data"))
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusCreated))
	w.WriteHeader(http.StatusCreated)
}

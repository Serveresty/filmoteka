package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

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
			logger.WarningLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status) + " | Details: " + string(err[:]))
		} else {
			logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status) + " | Details: " + string(err[:]))
		}
		w.WriteHeader(status)
		w.Write(err)
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status))
	w.WriteHeader(status)
}

func Registration(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

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
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(status) +
				" | Details: " + string(err[:]))

		w.WriteHeader(status)
		w.Write(err)
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
		logger.InfoLogger.Println(
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

	w.WriteHeader(http.StatusCreated)
}

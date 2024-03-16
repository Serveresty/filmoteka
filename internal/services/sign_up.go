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
		logger.InfoLogger.Println(r.Method + " " + r.URL.Path + " Status: " + http.StatusText(status) + " Details: " + string(err[:]))
		w.WriteHeader(status)
		w.Write(err)
		return
	}

	logger.InfoLogger.Println(r.Method + " " + r.URL.Path + " Status: " + http.StatusText(status))
	w.WriteHeader(status)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		w.WriteHeader(status)
		w.Write(err)
		return
	}

	err1 := json.NewDecoder(r.Body).Decode(&user)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error decoding json"))
		return
	}

	if !jwt.ValidateAuthData(&user) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error wrong input data"))
		return
	}

	hashPass, err1 := jwt.HashPassword(user.Password)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while hashing password"))
		return
	}

	user.Password = hashPass

	err1 = database.CreateNewUser(&user)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error user already created"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

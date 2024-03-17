package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"net/http"
)

// @Summary SignIn Page
// @Description This endpoint for sign-in page
// @Produce json
// @Success 200
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /sign-in [get]
func SignIn(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/sign-in" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /sign-in")
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

// @Summary Login
// @Description This endpoint for login
// @Accept json
// @Produce json
// @Success 200 {string} string "JSON с токеном"
// @Failure 400 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 403 {string} string "JSON с ошибками"
// @Failure 404 "Ничего"
// @Failure 405 "Ничего"
// @Failure 500	{string} string "JSON с ошибками"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/login" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /login")
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

	id, hashPass, err1 := database.GetUser(user.Email)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "user not found")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error user not found"))
		return
	}

	err1 = jwt.ComparePasswordHash(user.Password, hashPass)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "wrong password")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error wrong password"))
		return
	}

	roles, err1 := database.GetUserRoles(id)
	if err1 != nil {
		logger.WarningLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusInternalServerError) +
				" | Error: " + err1.Error() +
				" | Details: " + "no roles has been found")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error no roles has been found"))
		return
	}

	newToken, err1 := jwt.GenerateToken(id, roles)
	if err1 != nil {
		logger.WarningLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusInternalServerError) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while generating token")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while generating token"))
		return
	}

	var tokenMap []models.ErrorInfo
	tokenMap = append(tokenMap, models.ErrorInfo{Type: "token", Message: newToken})
	jsonResp, err1 := json.Marshal(tokenMap)
	if err1 != nil {
		logger.WarningLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusInternalServerError) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while marshal token")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marshal token"))
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK) + " | Token: " + newToken)
	r.Header.Set("Authorization", newToken)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

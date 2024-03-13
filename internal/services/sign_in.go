package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
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

	w.WriteHeader(status)
}

func Login(w http.ResponseWriter, r *http.Request) {
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

	id, hashPass, err1 := database.GetUser(user.Email)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error user not found"))
		return
	}

	err1 = jwt.ComparePasswordHash(user.Password, hashPass)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error wrong password"))
		return
	}

	roles, err1 := database.GetUserRoles(id)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error no roles has been found"))
		return
	}

	newToken, err1 := jwt.GenerateToken(id, roles)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while generating token"))
		return
	}

	tokenMap := make(map[string]string, 1)
	tokenMap["token"] = newToken
	jsonResp, err1 := json.Marshal(tokenMap)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal token"))
		return
	}

	r.Header.Set("Authorization", newToken)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResp))
}

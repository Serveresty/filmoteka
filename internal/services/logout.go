package services

import (
	"filmoteka/pkg/jwt"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
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
	r.Header.Del("Authorization")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("token removed"))
}

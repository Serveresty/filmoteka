package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/pkg/jwt"
	"net/http"
)

func GetActors(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/actors" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	_, err1 := jwt.ParseToken(token)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	actrs, err1 := database.GetActors()
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while getting actors"))
		return
	}

	jsonResp, err1 := json.Marshal(actrs)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal actors"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func AddNewActor(w http.ResponseWriter, r *http.Request) {

}

func EditInfoActor(w http.ResponseWriter, r *http.Request) {

}

func DeleteActor(w http.ResponseWriter, r *http.Request) {

}

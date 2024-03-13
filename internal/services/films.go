package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/pkg/jwt"
	"fmt"
	"net/http"
)

func GetFilms(w http.ResponseWriter, r *http.Request) {
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

	films, err1 := database.GetFilms("")
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while getting films"))
		fmt.Println(err1)
		return
	}

	jsonResp, err1 := json.Marshal(films)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal films"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResp))
}

func GetFilmsByFilter(w http.ResponseWriter, r *http.Request) {
	var filter string
	if r.URL.Path == "/films/search" {
		filter = r.URL.Query().Get("filter")
	} else {
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

	films, err1 := database.GetFilms(filter)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while getting films"))
		fmt.Println(err1)
		return
	}

	if films == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no films found by filter"))
		return
	}

	jsonResp, err1 := json.Marshal(films)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal films"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResp))
}

func AddNewFilm(w http.ResponseWriter, r *http.Request) {

}

func EditInfoFilm(w http.ResponseWriter, r *http.Request) {

}

func DeleteFilm(w http.ResponseWriter, r *http.Request) {

}

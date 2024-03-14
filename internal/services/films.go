package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/utils"
	"io"
	"net/http"
)

func GetFilms(w http.ResponseWriter, r *http.Request) {
	var sortBy string
	var direction string
	if r.URL.Path == "/films" {
		sortBy = r.URL.Query().Get("sort")
		direction = r.URL.Query().Get("direction")
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

	_, err1 := jwt.ParseToken(token)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	films, err1 := database.GetFilms("")
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while getting films"))
		return
	}

	films = utils.SortFilms(films, sortBy, direction)

	jsonResp, err1 := json.Marshal(films)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal films"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func GetFilmsByFilter(w http.ResponseWriter, r *http.Request) {
	var filter string
	var sortBy string
	var direction string
	if r.URL.Path == "/films/search" {
		filter = r.URL.Query().Get("filter")
		sortBy = r.URL.Query().Get("sort")
		direction = r.URL.Query().Get("direction")
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

	_, err1 := jwt.ParseToken(token)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	films, err1 := database.GetFilms(filter)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while getting films"))
		return
	}

	if films == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no films found by filter"))
		return
	}

	films = utils.SortFilms(films, sortBy, direction)

	jsonResp, err1 := json.Marshal(films)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal films"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func AddNewFilm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/new-film" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var filmToActor models.FilmToActor
	if err1 := json.Unmarshal(body, &filmToActor); err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.SetNewFilm(filmToActor)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
	}

	w.WriteHeader(http.StatusCreated)
}

func EditInfoFilm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit-film" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "PUT" {
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

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
}

func DeleteFilm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete-film" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "DELETE" {
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

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
}

package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"io"
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
	if r.URL.Path != "/new-actor" {
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

	var actorToFilm models.ActorToFilm
	if err1 := json.Unmarshal(body, &actorToFilm); err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.SetNewActor(actorToFilm)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func EditInfoActor(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit-actor" {
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

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actorToFilm models.ActorToFilm
	if err1 := json.Unmarshal(body, &actorToFilm); err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errr := database.EditActorInfo(actorToFilm)
	if errr != nil {
		errResp, err1 := json.Marshal(errr)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(errResp)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteActor(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete-actor" {
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

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actor models.Actor
	if err1 := json.Unmarshal(body, &actor); err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.DeleteActor(actor)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(errResp)
		return
	}

	w.WriteHeader(http.StatusOK)
}

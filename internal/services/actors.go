package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"fmt"
	"io"
	"net/http"
)

// @Summary GetActors
// @Description This endpoint for getting actors
// @Produce json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /actors [get]
func GetActors(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/actors" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /actors")
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
			logger.WarningLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status) + " | Details: " + string(err[:]))
		} else {
			logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status) + " | Details: " + string(err[:]))
		}

		w.WriteHeader(status)
		w.Write(err)
		return
	}

	_, err1 := jwt.ParseToken(token)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "not valid token")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	actrs, err1 := database.GetActors()
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while getting actors")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while getting actors"))
		return
	}

	jsonResp, err1 := json.Marshal(actrs)
	if err1 != nil {
		logger.WarningLogger.Println(
			r.Method + " | " +
				r.URL.Path + " | Status: " +
				http.StatusText(http.StatusInternalServerError) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while marhal films")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal actors"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// @Summary AddNewActor
// @Description This endpoint for adding new actor
// @Accept json
// @Produce json
// @Success 200
// @Success 201
// @Success 206
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /new-actor [post]
func AddNewActor(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/new-actor" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /new-actor")
		w.WriteHeader(http.StatusNotFound)
		return
	}

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
			logger.WarningLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status) + " | Details: " + string(err[:]))
		} else {
			logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(status) + " | Details: " + string(err[:]))
		}

		w.WriteHeader(status)
		w.Write(err)
		return
	}

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "not valid token")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusForbidden) +
				" | Details: " + "access denied")

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
	logger.InfoLogger.Println(
		r.Method + " | " +
			r.URL.Path +
			" | Status: " + http.StatusText(http.StatusOK) +
			" | Details: " + "access granted for userID: " + claims.Id)

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while reading request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actorToFilm models.ActorToFilm
	if err1 := json.Unmarshal(body, &actorToFilm); err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while parsing json")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.SetNewActor(actorToFilm)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			logger.WarningLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusInternalServerError) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while marshal array of errors")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusPartialContent) +
				" | Error: " + fmt.Sprintf("%v", errs) +
				" | Details: " + "request success with troubles")

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusCreated))
	w.WriteHeader(http.StatusCreated)
}

// @Summary EditInfoActor
// @Description This endpoint for edit actor's info
// @Accept json
// @Produce json
// @Success 200
// @Success 206
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /edit-actor [put]
func EditInfoActor(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/edit-actor" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /edit-actor")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "PUT" {
		logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected PUT for URL: " + r.URL.Path)
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

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "not valid token")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusForbidden) +
				" | Details: " + "access denied")

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
	logger.InfoLogger.Println(
		r.Method + " | " +
			r.URL.Path +
			" | Status: " + http.StatusText(http.StatusOK) +
			" | Details: " + "access granted for userID: " + claims.Id)

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while reading request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actorToFilm models.ActorToFilm
	if err1 := json.Unmarshal(body, &actorToFilm); err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while parsing json")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errr := database.EditActorInfo(actorToFilm)
	if errr != nil {
		errResp, err1 := json.Marshal(errr)
		if err1 != nil {
			logger.WarningLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusInternalServerError) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while marshal array of errors")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusPartialContent) +
				" | Error: " + fmt.Sprintf("%v", errr) +
				" | Details: " + "request success with troubles")

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

// @Summary DeleteActor
// @Description This endpoint for delete actor
// @Accept json
// @Produce json
// @Success 200
// @Success 206
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /delete-actor [delete]
func DeleteActor(w http.ResponseWriter, r *http.Request) {

	logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)

	if r.URL.Path != "/delete-actor" {
		logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /delete-actor")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "DELETE" {
		logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected DELETE for URL: " + r.URL.Path)
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

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "not valid token")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusForbidden) +
				" | Details: " + "access denied")

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
	logger.InfoLogger.Println(
		r.Method + " | " +
			r.URL.Path +
			" | Status: " + http.StatusText(http.StatusOK) +
			" | Details: " + "access granted for userID: " + claims.Id)

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while reading request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actor models.Actor
	if err1 := json.Unmarshal(body, &actor); err1 != nil {
		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusBadRequest) +
				" | Error: " + err1.Error() +
				" | Details: " + "error while parsing json")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.DeleteActor(actor)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			logger.WarningLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusInternalServerError) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while marshal array of errors")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		logger.InfoLogger.Println(
			r.Method + " | " + r.URL.Path +
				" | Status: " + http.StatusText(http.StatusPartialContent) +
				" | Error: " + fmt.Sprintf("%v", errs) +
				" | Details: " + "request success with troubles")

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

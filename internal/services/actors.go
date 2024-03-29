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
// @Security ApiKeyAuth
// @Tags actors
// @Description This endpoint for getting actors
// @Produce json
// @Success 200 {array} models.ActorToFilm "JSON с актёрами"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /actors [get]
func GetActors(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/actors" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /actors")
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected GET for URL: " + r.URL.Path)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		if !IsTesting {
			if status == http.StatusInternalServerError {
				logger.WarningLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			} else {
				logger.InfoLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			}
		}

		jsonResp, err := json.Marshal(err)
		if err != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}
		w.WriteHeader(status)
		w.Write(jsonResp)
		return
	}

	_, err1 := jwt.ParseToken(token)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "not valid token")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	actrs, err1 := database.GetActors()
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while getting actors")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while getting actors"))
		return
	}

	jsonResp, err1 := json.Marshal(actrs)
	if err1 != nil {
		if !IsTesting {
			logger.WarningLogger.Println(
				r.Method + " | " +
					r.URL.Path + " | Status: " +
					http.StatusText(http.StatusInternalServerError) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while marhal films")
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while marhal actors"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// @Summary AddNewActor
// @Security ApiKeyAuth
// @Tags actors
// @Description This endpoint for adding new actor
// @Accept json
// @Produce json
// @Param actor body models.Actor true "Данные фильма (используются все поля, кроме 'id')"
// @Success 201 "Ничего"
// @Success 206 {string} string "JSON с ошибками(если такие имеются)"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /new-actor [post]
func AddNewActor(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/new-actor" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /new-actor")
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected POST for URL: " + r.URL.Path)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		if !IsTesting {
			if status == http.StatusInternalServerError {
				logger.WarningLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			} else {
				logger.InfoLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			}
		}

		jsonResp, err := json.Marshal(err)
		if err != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}
		w.WriteHeader(status)
		w.Write(jsonResp)
		return
	}

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "not valid token")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusForbidden) +
					" | Details: " + "access denied")
		}

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
	if !IsTesting {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusOK) +
				" | Details: " + "access granted for userID: " + claims.Id)
	}

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while reading request body")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actorToFilm models.ActorToFilm
	if err1 := json.Unmarshal(body, &actorToFilm); err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while parsing json")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.SetNewActor(actorToFilm)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err1.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusPartialContent) +
					" | Error: " + fmt.Sprintf("%v", errs) +
					" | Details: " + "request success with troubles")
		}

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	if !IsTesting {
		logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusCreated))
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary EditInfoActor
// @Security ApiKeyAuth
// @Tags actors
// @Description This endpoint for edit actor's info
// @Accept json
// @Produce json
// @Param actor body models.Actor true "Данные фильма (используются все поля, кроме 'id')"
// @Success 200 "Ничего"
// @Success 206 {string} string "JSON с ошибками(если такие имеются)"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /edit-actor [put]
func EditInfoActor(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/edit-actor" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /edit-actor")
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "PUT" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected PUT for URL: " + r.URL.Path)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		if !IsTesting {
			if status == http.StatusInternalServerError {
				logger.WarningLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			} else {
				logger.InfoLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			}
		}

		jsonResp, err := json.Marshal(err)
		if err != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}
		w.WriteHeader(status)
		w.Write(jsonResp)
		return
	}

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "not valid token")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusForbidden) +
					" | Details: " + "access denied")
		}

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
	if !IsTesting {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusOK) +
				" | Details: " + "access granted for userID: " + claims.Id)
	}

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while reading request body")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actorToFilm models.ActorToFilm
	if err1 := json.Unmarshal(body, &actorToFilm); err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while parsing json")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errr := database.EditActorInfo(actorToFilm)
	if errr != nil {
		errResp, err1 := json.Marshal(errr)
		if err1 != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err1.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusPartialContent) +
					" | Error: " + fmt.Sprintf("%v", errr) +
					" | Details: " + "request success with troubles")
		}

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	if !IsTesting {
		logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK))
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary DeleteActor
// @Security ApiKeyAuth
// @Tags actors
// @Description This endpoint for delete actor
// @Accept json
// @Produce json
// @Param actor body models.Actor true "Данные фильма (используется только 'id')"
// @Success 200 "Ничего"
// @Success 206 {string} string "JSON с ошибками(если такие имеются)"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /delete-actor [delete]
func DeleteActor(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/delete-actor" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /delete-actor")
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "DELETE" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request method: " + r.Method + ", expected DELETE for URL: " + r.URL.Path)
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	path := r.URL.Path
	status, err := jwt.CheckAuthorization(token, path)
	if err != nil {
		if !IsTesting {
			if status == http.StatusInternalServerError {
				logger.WarningLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			} else {
				logger.InfoLogger.Printf(r.Method+" | "+r.URL.Path+" | Status: "+http.StatusText(status)+" | Details: %v", err)
			}
		}

		jsonResp, err := json.Marshal(err)
		if err != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}
		w.WriteHeader(status)
		w.Write(jsonResp)
		return
	}

	claims, err1 := jwt.ParseToken(token)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "not valid token")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error not valid token"))
		return
	}

	if !jwt.HasUserAccess(*claims) {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusForbidden) +
					" | Details: " + "access denied")
		}

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("error access denied"))
		return
	}
	if !IsTesting {
		logger.InfoLogger.Println(
			r.Method + " | " +
				r.URL.Path +
				" | Status: " + http.StatusText(http.StatusOK) +
				" | Details: " + "access granted for userID: " + claims.Id)
	}

	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while reading request body")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while reading request body"))
		return
	}

	var actor models.Actor
	if err1 := json.Unmarshal(body, &actor); err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while parsing json")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while parsing json"))
		return
	}

	errs := database.DeleteActor(actor)
	if errs != nil {
		errResp, err1 := json.Marshal(errs)
		if err1 != nil {
			if !IsTesting {
				logger.WarningLogger.Println(
					r.Method + " | " + r.URL.Path +
						" | Status: " + http.StatusText(http.StatusInternalServerError) +
						" | Error: " + err1.Error() +
						" | Details: " + "error while marshal array of errors")
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error while marhal errs"))
			return
		}

		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " + r.URL.Path +
					" | Status: " + http.StatusText(http.StatusPartialContent) +
					" | Error: " + fmt.Sprintf("%v", errs) +
					" | Details: " + "request success with troubles")
		}

		w.WriteHeader(http.StatusPartialContent)
		w.Write(errResp)
		return
	}

	if !IsTesting {
		logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK))
	}
	w.WriteHeader(http.StatusOK)
}

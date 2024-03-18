package services

import (
	"encoding/json"
	"filmoteka/internal/database"
	"filmoteka/internal/models"
	"filmoteka/pkg/jwt"
	"filmoteka/pkg/logger"
	"filmoteka/pkg/utils"
	"fmt"
	"io"
	"net/http"
)

// @Summary GetFilms
// @Security ApiKeyAuth
// @Tags films
// @Description This endpoint for getting films
// @Produce json
// @Success 200 {array} models.FilmToActor "JSON с фильмами"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /films [get]
func GetFilms(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	var sortBy string
	var direction string
	if r.URL.Path == "/films" {
		sortBy = r.URL.Query().Get("sort")
		direction = r.URL.Query().Get("direction")
	} else {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /films")
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

	films, err1 := database.GetFilms("")
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while getting films")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while getting films"))
		return
	}

	films = utils.SortFilms(films, sortBy, direction)

	jsonResp, err1 := json.Marshal(films)
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
		w.Write([]byte("error while marhal films"))
		return
	}

	if !IsTesting {
		logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// @Summary GetFilmsByFilter
// @Security ApiKeyAuth
// @Tags films
// @Description This endpoint for getting films by filter
// @Produce json
// @Success 200 {array} models.FilmToActor "JSON с фильмами"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /films/search [get]
func GetFilmsByFilter(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	var filter string
	var sortBy string
	var direction string
	if r.URL.Path == "/films/search" {
		filter = r.URL.Query().Get("filter")
		sortBy = r.URL.Query().Get("sort")
		direction = r.URL.Query().Get("direction")
	} else {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /films/search")
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

	films, err1 := database.GetFilms(filter)
	if err1 != nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusBadRequest) +
					" | Error: " + err1.Error() +
					" | Details: " + "error while getting films")
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error while getting films"))
		return
	}

	if films == nil {
		if !IsTesting {
			logger.InfoLogger.Println(
				r.Method + " | " +
					r.URL.Path +
					" | Status: " + http.StatusText(http.StatusOK) +
					" | Details: " + "no films found by filter")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no films found by filter"))
		return
	}

	films = utils.SortFilms(films, sortBy, direction)

	jsonResp, err1 := json.Marshal(films)
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
		w.Write([]byte("error while marhal films"))
		return
	}

	logger.InfoLogger.Println(r.Method + " | " + r.URL.Path + " | Status: " + http.StatusText(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// @Summary AddNewFilm
// @Security ApiKeyAuth
// @Tags films
// @Description This endpoint for adding new film
// @Accept json
// @Produce json
// @Param film body models.Film true "Данные фильма (используются все поля, кроме 'id')"
// @Success 201 "Ничего"
// @Success 206 {array} models.ErrorInfo "JSON с ошибками(если такие имеются)"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /new-film [post]
func AddNewFilm(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/new-film" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /new-film")
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

	var filmToActor models.FilmToActor
	if err1 := json.Unmarshal(body, &filmToActor); err1 != nil {
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

	errs := database.SetNewFilm(filmToActor)
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
			w.Write([]byte("error while marhal array of errors"))
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

// @Summary EditInfoFilm
// @Security ApiKeyAuth
// @Tags films
// @Description This endpoint for edit film's info
// @Accept json
// @Produce json
// @Param film body models.Film true "Данные фильма (используются все поля, кроме 'id')"
// @Success 200 "Ничего"
// @Success 206 {string} string "JSON с ошибками(если такие имеются)"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /edit-film [put]
func EditInfoFilm(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/edit-film" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /edit-film")
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

	var filmToActor models.FilmToActor
	if err1 := json.Unmarshal(body, &filmToActor); err1 != nil {
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

	errr := database.EditFilmInfo(filmToActor)
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

// @Summary DeleteFilm
// @Security ApiKeyAuth
// @Tags films
// @Description This endpoint for delete film
// @Accept json
// @Produce json
// @Param film body models.Film true "Данные фильма (используется 'id')"
// @Success 200 "Ничего"
// @Success 206 {string} string "JSON с ошибками(если такие имеются)"
// @Failure 400,401,403 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Failure 404,405 "Ничего"
// @Failure 500 {string} string "JSON с ошибками, либо строка(в зависимости от возвращающего метода)"
// @Router /delete-film [delete]
func DeleteFilm(w http.ResponseWriter, r *http.Request) {

	if !IsTesting {
		logger.InfoLogger.Println("Handling " + r.Method + " request for: " + r.URL.Path)
	}

	if r.URL.Path != "/delete-film" {
		if !IsTesting {
			logger.InfoLogger.Println("Invalid request URL: " + r.URL.Path + ", expected URL: /edit-film")
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

	var film models.Film
	if err1 := json.Unmarshal(body, &film); err1 != nil {
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

	errs := database.DeleteFilm(film)
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

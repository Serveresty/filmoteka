package database

import (
	"database/sql"
	"filmoteka/internal/models"
)

func SetNewActor(actorToFilm models.ActorToFilm) []models.ErrorInfo {
	var actor models.Actor
	var errorsPull []models.ErrorInfo
	row := DB.QueryRow(`
		SELECT * FROM actors WHERE name=$1 AND gender=$2 AND birth_date=$3`,
		actorToFilm.Actors.Name, actorToFilm.Actors.Gender, actorToFilm.Actors.BirthDate)

	err := row.Scan(&actor)
	if err != sql.ErrNoRows {
		errorsPull = append(errorsPull, models.ErrorInfo{Message: "actor already has in table"})
		return errorsPull
	}

	var actorId int
	err = DB.QueryRow(`INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING actor_id`, actorToFilm.Actors.Name, actorToFilm.Actors.Gender, actorToFilm.Actors.BirthDate).Scan(&actorId)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		return errorsPull
	}

	var filmsIds []int
	for _, val := range actorToFilm.Movies {
		var film models.Film
		row := DB.QueryRow(`
			SELECT * FROM films WHERE name=$1 AND description=$2 AND release_date=$3 AND rate=$4`,
			val.Title, val.Description, val.ReleaseDate, val.Rate)

		err := row.Scan(&film)
		if err != sql.ErrNoRows {
			filmsIds = append(filmsIds, film.ID)
			continue
		}

		var filmId int
		err = DB.QueryRow(`
			INSERT INTO films (name, description, release_date, rate) VALUES ($1, $2, $3) RETURNING film_id`,
			val.Title, val.Description, val.ReleaseDate, val.Rate).Scan(&filmId)
		if err != nil {
			errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
			return errorsPull
		}
		filmsIds = append(filmsIds, filmId)
	}

	for _, flm := range filmsIds {
		_, err = DB.Exec(`INSERT INTO actors_films (actor_id, film_id) VALUES ($1, $2)`, flm, actorId)
		if err != nil {
			errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		}
	}

	if len(errorsPull) > 0 {
		return errorsPull
	}

	return nil
}

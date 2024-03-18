package database

import (
	"database/sql"
	"filmoteka/internal/models"
)

func SetNewFilm(filmToActor models.FilmToActor) []models.ErrorInfo {
	var film models.Film
	var errorsPull []models.ErrorInfo
	row := DB.QueryRow(`
		SELECT * FROM films WHERE name=$1 AND description=$2 AND release_date=$3 AND rate=$4`,
		filmToActor.Movies.Title, filmToActor.Movies.Description, filmToActor.Movies.ReleaseDate, filmToActor.Movies.Rate)

	err := row.Scan(&film)
	if err != sql.ErrNoRows {
		errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "film already has in table"})
		return errorsPull
	}

	var filmId int
	err = DB.QueryRow(`INSERT INTO films (name, description, release_date, rate) VALUES ($1, $2, $3, $4) RETURNING film_id`, filmToActor.Movies.Title, filmToActor.Movies.Description, filmToActor.Movies.ReleaseDate, filmToActor.Movies.Rate).Scan(&filmId)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while insert film"})
		return errorsPull
	}

	var actorsIds []int
	for _, val := range filmToActor.Actors {
		var actor models.Actor
		row := DB.QueryRow(`
			SELECT * FROM actors WHERE name=$1 AND gender=$2 AND birth_date=$3`,
			val.Name, val.Gender, val.BirthDate)

		err := row.Scan(&actor)
		if err != sql.ErrNoRows {
			actorsIds = append(actorsIds, actor.ID)
			continue
		}

		var actorId int
		err = DB.QueryRow(`
			INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING actor_id`,
			val.Name, val.Gender, val.BirthDate).Scan(&actorId)
		if err != nil {
			errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while insert actor"})
			return errorsPull
		}
		actorsIds = append(actorsIds, actorId)
	}

	for _, act := range actorsIds {
		_, err = DB.Exec(`INSERT INTO actors_films (actor_id, film_id) VALUES ($1, $2)`, act, filmId)
		if err != nil {
			errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while insert actors films"})
		}
	}

	if len(errorsPull) > 0 {
		return errorsPull
	}

	return nil
}

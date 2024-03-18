package database

import (
	"filmoteka/internal/models"
)

func EditFilmInfo(filmToActor models.FilmToActor) []models.ErrorInfo {
	var errorsPull []models.ErrorInfo
	_, err := DB.Exec(`
		UPDATE films SET name=$1, description=$2, release_date=$3, rate=$4 
		WHERE film_id=$5;`,
		filmToActor.Movies.Title, filmToActor.Movies.Description, filmToActor.Movies.ReleaseDate, filmToActor.Movies.Rate, filmToActor.Movies.ID)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while update film's info"})
		return errorsPull
	}

	_, err = DB.Exec(`DELETE FROM actors_films WHERE film_id=$1`, filmToActor.Movies.ID)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while update film's info on delete"})
		return errorsPull
	}

	for _, act := range filmToActor.Actors {
		_, err = DB.Exec(`INSERT INTO actors_films (actor_id, film_id) VALUE($1, $2)`, act.ID, filmToActor.Movies.ID)
		if err != nil {
			errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while update film's info on insert"})
		}
	}

	return nil
}

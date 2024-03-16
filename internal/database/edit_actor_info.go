package database

import "filmoteka/internal/models"

func EditActorInfo(actorToFilm models.ActorToFilm) []models.ErrorInfo {
	var errorsPull []models.ErrorInfo
	_, err := DB.Exec(`
		UPDATE actors SET name=$1, gender=$2, birth_date=$3 
		WHERE actor_id=$4;`,
		actorToFilm.Actors.Name, actorToFilm.Actors.Gender, actorToFilm.Actors.BirthDate, actorToFilm.Actors.ID)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		return errorsPull
	}

	_, err = DB.Exec(`DELETE FROM actors_films WHERE actor_id=$1`, actorToFilm.Actors.ID)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		return errorsPull
	}

	for _, flm := range actorToFilm.Movies {
		_, err = DB.Exec(`INSERT INTO actors_films (actor_id, film_id) VALUE($1, $2)`, actorToFilm.Actors.ID, flm.ID)
		if err != nil {
			errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		}
	}

	return nil
}

package database

import (
	"filmoteka/internal/models"
)

func GetActors() ([]models.ActorToFilm, error) {
	rows, err := DB.Query(`
	SELECT actors.*, films.* 
	FROM actors 
	LEFT JOIN actors_films ON actors.actor_id = actors_films.actor_id 
	LEFT JOIN films ON actors_films.film_id = films.film_id
	WHERE films.film_id IS NOT NULL
	
	UNION ALL
	
	SELECT actors.*, 0 AS film_id, '' AS name, '' AS description, '0001-01-01' AS release_date, 0.0 AS rate 
	FROM actors 
	LEFT JOIN actors_films ON actors.actor_id = actors_films.actor_id 
	LEFT JOIN films ON actors_films.film_id = films.film_id
	WHERE actors.actor_id NOT IN (SELECT actor_id FROM actors_films);`)
	if err != nil {
		return nil, err
	}

	var actors = make(map[models.Actor][]models.Film)
	for rows.Next() {
		var film models.Film
		var actor models.Actor
		var gender string
		err = rows.Scan(&actor.ID, &actor.Name, &gender, &actor.BirthDate, &film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rate)
		if err != nil {
			return nil, err
		}

		actor.Gender = gender

		if val, ok := actors[actor]; ok {
			val = append(val, film)
			actors[actor] = val
		} else {
			actors[actor] = []models.Film{film}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []models.ActorToFilm
	for act, fil := range actors {
		var res models.ActorToFilm
		res.Actors = act
		res.Movies = fil
		result = append(result, res)
	}
	return result, nil
}

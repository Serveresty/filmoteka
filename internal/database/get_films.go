package database

import (
	"database/sql"
	"filmoteka/internal/models"
)

func GetFilms(filter string) ([]models.Film, error) {
	var query string = `
		SELECT DISTINCT f.*, a.*
		FROM films f
		JOIN actors_films af ON f.film_id = af.film_id 
		JOIN actors a ON af.actor_id = a.actor_id
		`

	if filter != "" {
		query += ` WHERE LOWER(f.name) LIKE LOWER($1) OR LOWER(a.name) LIKE LOWER($1)`
	}
	query += `;`

	var rows *sql.Rows
	var err error

	if filter != "" {
		rows, err = DB.Query(query, filter)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = DB.Query(query)
		if err != nil {
			return nil, err
		}
	}

	var filmsMap = make(map[*models.Film][]models.Actor)
	var films []models.Film
	for rows.Next() {
		var film models.Film
		var actors models.Actor
		var gender string
		err = rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rate, &actors.ID, &actors.Name, &gender, &actors.BirthDate)
		if err != nil {
			return nil, err
		}

		actors.Gender = []rune(gender)[0]

		if _, ok := filmsMap[&film]; !ok {
			film.Actors = append(film.Actors, actors)
			filmsMap[&film] = film.Actors
		} else {
			filmsMap[&film] = append(filmsMap[&film], actors)
		}

		// film.Actors = append(film.Actors, actors)
		// films = append(films, film)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for film, actors := range filmsMap {
		film.Actors = actors
		films = append(films, *film)
	}

	return films, nil
}

/*
`
	SELECT DISTINCT f.*
	FROM films f
	JOIN actors_films af ON f.film_id = af.film_id
	JOIN actors a ON af.actor_id = a.actor_id
	WHERE LOWER(f.name) LIKE LOWER('film name') OR LOWER(a.name) LIKE LOWER('actor name');`
*/

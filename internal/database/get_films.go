package database

import (
	"database/sql"
	"filmoteka/internal/models"
	"fmt"
)

func GetFilms(filter string) ([]models.FilmToActor, error) {
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

	var films = make(map[models.Film][]models.Actor)
	for rows.Next() {
		var film models.Film
		var actors models.Actor
		var gender string
		err = rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rate, &actors.ID, &actors.Name, &gender, &actors.BirthDate)
		if err != nil {
			return nil, err
		}

		actors.Gender = []rune(gender)[0]

		if val, ok := films[film]; ok {
			val = append(val, actors)
			films[film] = val
		} else {
			films[film] = []models.Actor{actors}
		}

		//films[&film] = append(films[&film], actors)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(films)

	var result []models.FilmToActor
	for fil, act := range films {
		var res models.FilmToActor
		res.Movies = fil
		res.Actors = act
		result = append(result, res)
	}

	return result, nil
}

/*
`
	SELECT DISTINCT f.*
	FROM films f
	JOIN actors_films af ON f.film_id = af.film_id
	JOIN actors a ON af.actor_id = a.actor_id
	WHERE LOWER(f.name) LIKE LOWER('film name') OR LOWER(a.name) LIKE LOWER('actor name');`
*/

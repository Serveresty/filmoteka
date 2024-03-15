package database

import "filmoteka/internal/models"

func DeleteFilm(film models.Film) []models.ErrorInfo {
	var errorsPull []models.ErrorInfo
	_, err := DB.Exec(`DELETE FROM films WHERE film_id=$1`, film.ID)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		return errorsPull
	}
	return nil
}

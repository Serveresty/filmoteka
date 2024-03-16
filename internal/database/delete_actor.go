package database

import "filmoteka/internal/models"

func DeleteActor(actor models.Actor) []models.ErrorInfo {
	var errorsPull []models.ErrorInfo
	_, err := DB.Exec(`DELETE FROM actors WHERE actor_id=$1`, actor.ID)
	if err != nil {
		errorsPull = append(errorsPull, models.ErrorInfo{Message: err.Error()})
		return errorsPull
	}
	return nil
}

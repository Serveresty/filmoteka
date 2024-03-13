package utils

import (
	"filmoteka/internal/models"
	"sort"
)

func SortFilms(filmsAndActors []models.FilmToActor, sortBy string, direction string) []models.FilmToActor {
	if sortBy == "" && direction != "" {
		sortBy = "rate"
	}
	switch sortBy {
	case "title":
		if direction == "asc" {
			sort.SliceStable(filmsAndActors, func(i, j int) bool {
				return filmsAndActors[i].Movies.Title > filmsAndActors[j].Movies.Title
			})
		}
		if direction == "desc" {
			sort.SliceStable(filmsAndActors, func(i, j int) bool {
				return filmsAndActors[i].Movies.Title < filmsAndActors[j].Movies.Title
			})
		}
	case "release-date":
		if direction == "asc" {
			sort.SliceStable(filmsAndActors, func(i, j int) bool {
				return filmsAndActors[i].Movies.ReleaseDate.Before(filmsAndActors[j].Movies.ReleaseDate)
			})
		}
		if direction == "desc" {
			sort.SliceStable(filmsAndActors, func(i, j int) bool {
				return filmsAndActors[i].Movies.ReleaseDate.After(filmsAndActors[j].Movies.ReleaseDate)
			})
		}
	case "rate":
		if direction == "asc" {
			sort.SliceStable(filmsAndActors, func(i, j int) bool {
				return filmsAndActors[i].Movies.Rate < filmsAndActors[j].Movies.Rate
			})
		}
		if direction == "desc" {
			sort.SliceStable(filmsAndActors, func(i, j int) bool {
				return filmsAndActors[i].Movies.Rate > filmsAndActors[j].Movies.Rate
			})
		}
	default:
		sort.SliceStable(filmsAndActors, func(i, j int) bool {
			return filmsAndActors[i].Movies.Rate > filmsAndActors[j].Movies.Rate
		})
	}

	return filmsAndActors
}

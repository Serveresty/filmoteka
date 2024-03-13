package models

type ActorToFilm struct {
	Actors Actor  `json:"actors"`
	Movies []Film `json:"movies"`
}

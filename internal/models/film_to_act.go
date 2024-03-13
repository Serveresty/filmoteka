package models

type FilmToActor struct {
	Movies Film    `json:"movies"`
	Actors []Actor `json:"actors"`
}

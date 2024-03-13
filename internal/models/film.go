package models

import "time"

type Film struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rate        float64   `json:"rate"`
}

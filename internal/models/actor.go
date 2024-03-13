package models

import "time"

type Actor struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

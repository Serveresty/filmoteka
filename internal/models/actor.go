package models

import "time"

type Actor struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Gender    rune      `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

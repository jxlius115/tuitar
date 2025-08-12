// internal/models/tab.go
package models

import (
	// "encoding/json"
	"time"
)

type Tab struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Artist        string    `json:"artist" db:"artist"`
	Content       [6]string `json:"content" db:"content"` // 6 strings
	Tuning        [6]string `json:"tuning" db:"tuning"`   // E A D G B e
	Tempo         int       `json:"tempo" db:"tempo"`
	TimeSignature string    `json:"time_signature" db:"time_signature"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

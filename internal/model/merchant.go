package model

import (
	"database/sql"
	"time"
)

type Merchant struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Category  sql.NullString `json:"category" db:"category"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
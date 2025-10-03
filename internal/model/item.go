package model

import (
	"database/sql"
	"time"
)

type Item struct {
	ID         int64          `json:"id" db:"id"`
	MerchantID int64          `json:"merchant_id" db:"merchant_id"`
	Name       string         `json:"name" db:"name"`
	Category   sql.NullString `json:"category" db:"category"`
	Price      sql.NullInt64  `json:"price" db:"price"`
	FileURI    sql.NullString `json:"file_uri" db:"file_uri"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}
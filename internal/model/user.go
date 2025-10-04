package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Username     sql.NullString
	Email        sql.NullString
	PasswordHash string
	Role         int16
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

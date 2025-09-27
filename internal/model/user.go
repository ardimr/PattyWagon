package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int64
	Email             sql.NullString
	Phone             sql.NullString
	PasswordHash      string
	Name              sql.NullString
	FileID            sql.NullInt64
	FileURI           sql.NullString
	FileThumbnailURI  sql.NullString
	BankAccountName   sql.NullString
	BankAccountHolder sql.NullString
	BankAccountNumber sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

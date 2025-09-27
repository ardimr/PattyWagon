package model

import "database/sql"

type File struct {
	ID           int64        `db:"id"`
	Uri          string       `db:"uri"`
	ThumbnailUri string       `db:"thumbnail_uri"`
	SizeInBytes  int64        `db:"size_in_bytes"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}

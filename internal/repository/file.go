package repository

import (
	"WeenieHut/internal/model"
	"WeenieHut/observability"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (q *Queries) InsertFile(ctx context.Context, data model.File) (res model.File, err error) {
	ctx, span := observability.Tracer.Start(ctx, "repository.insert_file")
	defer span.End()

	query := `
		INSERT INTO files (
			uri, thumbnail_uri, created_at, updated_at
		) VALUES (
			$1, $2, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`

	err = q.db.QueryRowContext(ctx, query,
		data.Uri,
		data.ThumbnailUri,
	).Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		return model.File{}, fmt.Errorf("error inserting file: %w", err)
	}

	return data, nil
}

func (q *Queries) GetFileUpload(ctx context.Context, id int64) (model.File, error) {
	query := `
		SELECT id, uri, thumbnail_uri, created_at, updated_at FROM files
		WHERE id = $1
	`

	row := q.db.QueryRowContext(ctx, query, id)
	var f model.File
	if err := row.Scan(
		&f.ID,
		&f.Uri,
		&f.ThumbnailUri,
		&f.CreatedAt,
		&f.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.File{}, errors.New("file not found")
		}
		return model.File{}, err
	}

	return f, nil
}

func (q *Queries) FileExists(ctx context.Context, fileID string) (bool, error) {
	query := `
		SELECT 1 FROM files
		WHERE id = $1
	`

	var exists int
	err := q.db.QueryRowContext(ctx, query, fileID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (q *Queries) GetFileByFileID(ctx context.Context, fileID string) (res model.File, err error) {
	query := `SELECT id, uri, thumbnail_uri, created_at, updated_at FROM files WHERE id = $1`
	err = q.db.QueryRowContext(ctx, query, fileID).Scan(&res.ID, &res.Uri, &res.ThumbnailUri, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return model.File{}, err
	}
	return res, nil
}

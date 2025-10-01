package repository

import (
	"PattyWagon/internal/model"
	"context"
	"fmt"
)

func (q *Queries) InsertMerchant(ctx context.Context, data model.Merchant) (res int64, err error) {
	query := `
		INSERT INTO merchants (
			user_id, name, category, image_url, latitude, longitude, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(), NOW()
		)
		RETURNING id
	`

	err = q.db.QueryRowContext(ctx, query,
		data.UserID,
		data.Name,
		data.Category,
		data.ImageURL,
		data.Latitude,
		data.Longitude,
	).Scan(&data.ID)

	if err != nil {
		return 0, fmt.Errorf("error inserting file: %w", err)
	}

	return
}

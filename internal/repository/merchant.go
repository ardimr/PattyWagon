package repository

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"database/sql"
	"errors"
)

const (
	selectMerchantByIdQuery = `SELECT id, user_id, name, category, image_url, latitude, longitude, created_at, updated_at FROM merchants WHERE id = $1`
)

func (q *Queries) GetMerchantByID(ctx context.Context, id int64) (model.Merchant, error) {
	query := selectMerchantByIdQuery
	row := q.db.QueryRowContext(ctx, query, id)
	var m model.Merchant
	if err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.Name,
		&m.Category,
		&m.ImageURL,
		&m.Latitude,
		&m.Longitude,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Merchant{}, constants.ErrMerchantNotFound
		}
		return model.Merchant{}, err
	}

	return m, nil
}
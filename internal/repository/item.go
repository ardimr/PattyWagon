package repository

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"database/sql"
	"errors"
)

const (
	selectItemByIdQuery = `SELECT id, merchant_id, name, category, price, image_url, created_at, updated_at FROM items WHERE id = $1`
)

func (q *Queries) GetItemByID(ctx context.Context, id int64) (model.Item, error) {
	query := selectItemByIdQuery
	row := q.db.QueryRowContext(ctx, query, id)
	var i model.Item
	if err := row.Scan(
		&i.ID,
		&i.MerchantID,
		&i.Name,
		&i.Category,
		&i.Price,
		&i.ImageURL,
		&i.CreatedAt,
		&i.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Item{}, constants.ErrItemNotFound
		}
		return model.Item{}, err
	}

	return i, nil
}

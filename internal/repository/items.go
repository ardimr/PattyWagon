package repository

import (
	"PattyWagon/internal/model"
	"context"
)

func (q *Queries) CreateItems(ctx context.Context, item model.Item) (int64, error) {
	query := `
		INSERT INTO items (merchant_id, name, category, price, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id int64
	err := q.db.QueryRowContext(ctx, query,
		item.MerchantID,
		item.Name,
		item.Category,
		item.Price,
		item.ImageURL,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

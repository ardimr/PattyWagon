package repository

import (
	"PattyWagon/internal/model"
	"context"
	"fmt"
	"strings"
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

func (q *Queries) GetItems(ctx context.Context, filter model.FilterItem) (res []model.Item, err error) {
	query := `
		SELECT id, name, category, image_url, latitude, longitude, created_at
		FROM merchants
	`
	conds := []string{}
	args := []interface{}{}
	argIdx := 1

	// filter by MerchantID
	if filter.ItemID != 0 {
		conds = append(conds, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, filter.ItemID)
		argIdx++
	}

	// filter by Name
	if filter.Name != "" {
		conds = append(conds, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Name+"%")
		argIdx++
	}

	// filter by category
	if filter.ProductCategory != "" {
		conds = append(conds, fmt.Sprintf("LOWER(category) = LOWER($%d)", argIdx))
		args = append(args, filter.ProductCategory)
		argIdx++
	}

	// conditions
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	// sorting createdAt
	if strings.ToLower(filter.CreatedAt) == "asc" {
		query += " ORDER BY created_at ASC"
	} else if strings.ToLower(filter.CreatedAt) == "desc" {
		query += " ORDER BY created_at DESC"
	} else {
		query += " ORDER BY created_at DESC"
	}

	// pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 5
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var m model.Item
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Category,
			&m.Price,
			&m.ImageURL,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	res = items
	return
}

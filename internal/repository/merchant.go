package repository

import (
	"PattyWagon/internal/model"
	"context"
	"fmt"
	"strings"
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

func (q *Queries) GetMerchants(ctx context.Context, filter model.FilterMerchant) (res []model.Merchant, err error) {
	query := `
		SELECT id, name, category, image_url, latitude, longitude, created_at
		FROM merchants
	`
	conds := []string{}
	args := []interface{}{}
	argIdx := 1

	// filter by MerchantID
	if filter.MerchantID != 0 {
		conds = append(conds, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, filter.MerchantID)
		argIdx++
	}

	// filter by Name
	if filter.Name != "" {
		conds = append(conds, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+filter.Name+"%")
		argIdx++
	}

	// filter by category
	if filter.MerchantCategory != "" {
		conds = append(conds, fmt.Sprintf("LOWER(category) = LOWER($%d)", argIdx))
		args = append(args, filter.MerchantCategory)
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

	var merchants []model.Merchant
	for rows.Next() {
		var m model.Merchant
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Category,
			&m.ImageURL,
			&m.Latitude,
			&m.Longitude,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		merchants = append(merchants, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	res = merchants
	return
}

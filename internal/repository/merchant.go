package repository

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (q *Queries) InsertMerchant(ctx context.Context, data model.Merchant) (res int64, err error) {
	query := `
		INSERT INTO merchants (
			user_id, name, category, image_url, latitude, longitude, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW()
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

	return data.ID, nil
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

func (q *Queries) MerchantExists(ctx context.Context, merchantID int64) (res bool, err error) {
	var exists bool
	err = q.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM merchants WHERE id=$1)", merchantID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, constants.ErrMerchantNotFound
	}

	return exists, nil
}

func (q *Queries) BulkInsertMerchantLocations(ctx context.Context, locations []model.MerchantLocation) error {
	if len(locations) == 0 {
		return nil
	}

	query := `INSERT INTO merchant_locations (merchant_id, h3_index, resolution, created_at, updated_at) VALUES `

	values := []interface{}{}
	placeholders := []string{}

	for i, loc := range locations {
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, NOW(), NOW())", i*3+1, i*3+2, i*3+3))
		values = append(values, loc.MerchantID, loc.H3Index, loc.Resolution)
	}

	query += strings.Join(placeholders, ", ")

	_, err := q.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("error bulk inserting merchant locations: %w", err)
	}

	return nil
}

func (q *Queries) GetMerchantCount(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM merchants`
	err := q.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

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

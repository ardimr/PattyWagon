package repository

import (
	"PattyWagon/internal/model"
	"PattyWagon/logger"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	baseMerchantQuery = `
WITH matching_merchants AS (
    -- Base merchants query
    SELECT DISTINCT ml.merchant_id
    FROM merchant_locations ml
    INNER JOIN merchants m ON m.id = ml.merchant_id`

	unionItemQuery = `
    UNION
    -- Merchants matching by item name
    SELECT DISTINCT ml.merchant_id
    FROM merchant_locations ml
    INNER JOIN merchants m ON m.id = ml.merchant_id
    INNER JOIN items i ON i.merchant_id = m.id
  	`

	finalSelectQuery = `
)
SELECT
  m.id,
  m.name,
  m.category,
  m.image_url,
  m.latitude,
  m.longitude,
  m.created_at,
  json_agg(
    json_build_object(
      'id', i.id,
      'name', i.name,
      'product_category', i.category,
      'price', i.price,
      'image_url', i.image_url,
      'created_at', i.created_at AT TIME ZONE 'UTC'
    )
  ) as items
FROM matching_merchants
INNER JOIN merchants as m on m.id = matching_merchants.merchant_id
INNER JOIN items as i on i.merchant_id = matching_merchants.merchant_id
GROUP BY m.id`
)

const (
	getMerchantWithItems = `
SELECT
  m.id,
  m.name,
  m.category,
  m.image_url,
  m.latitude,
  m.longitude,
  m.created_at,
  json_agg(
    json_build_object(
      'id', i.id,
      'name', i.name,
      'product_category', i.category,
      'price', i.price,
      'image_url', i.image_url,
      'created_at', i.created_at AT TIME ZONE 'UTC'
    )
  ) as items
FROM merchants as m
LEFT JOIN items as i on i.merchant_id = m.id
WHERE m.id=$1
GROUP BY m.id`
)

func (q *Queries) GetMerchantByID(ctx context.Context, id int64) (model.Merchant, error) {
	panic("unimplemented")
}

func (q *Queries) GetMerchantByCellID(ctx context.Context, cellID int64) (model.Merchant, error) {
	panic("unimplemented")
}

func (q *Queries) ListMerchantWithItems(ctx context.Context, filter model.ListMerchantWithItemParams) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)
	// Build the query dynamically
	query := baseMerchantQuery
	args := []interface{}{}
	argIdx := 1

	// Add conditions to base query
	baseConds := []string{}

	if filter.Cell != nil {
		baseConds = append(baseConds, fmt.Sprintf(" ml.h3_index = $%d", argIdx))
		args = append(args, filter.Cell.CellID)
		argIdx++
	}

	if filter.MerchantCategory != nil {
		baseConds = append(baseConds, fmt.Sprintf(" m.category = $%d", argIdx))
		args = append(args, *filter.MerchantCategory)
		argIdx++
	}

	if filter.Name != nil {
		baseConds = append(baseConds, fmt.Sprintf(" m.name ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Name+"%")
		argIdx++
	}

	// Add base conditions to query
	if len(baseConds) > 0 {
		query += " " + "WHERE" + " " + strings.Join(baseConds, " AND")
	}

	// Add UNION conditionally if Name filter exists
	if filter.Name != nil {
		// Add same conditions for union query
		unionConds := []string{}
		query += unionItemQuery

		// Add h3_index for union query
		if filter.Cell != nil {
			unionConds = append(unionConds, fmt.Sprintf(" ml.h3_index = $%d", argIdx))
			args = append(args, filter.Cell.CellID)
			argIdx++
		}

		if filter.MerchantCategory != nil {
			unionConds = append(unionConds, fmt.Sprintf(" m.category = $%d", argIdx))
			args = append(args, *filter.MerchantCategory)
			argIdx++
		}

		unionConds = append(unionConds, fmt.Sprintf(" i.name ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Name+"%")
		argIdx++

		if len(unionConds) > 0 {
			query += " " + "WHERE " + strings.Join(unionConds, " AND")
		}
	}

	// Complete the query
	query += finalSelectQuery

	// Add final WHERE conditions if needed
	finalConds := []string{}
	if filter.MerchantID != nil {
		finalConds = append(finalConds, fmt.Sprintf("WHERE m.id = $%d", argIdx))
		args = append(args, *filter.MerchantID)
		argIdx++
	}

	if len(finalConds) > 0 {
		// Insert WHERE clause before GROUP BY
		query = strings.Replace(query, "GROUP BY m.id", strings.Join(finalConds, " AND ")+" GROUP BY m.id", 1)
	}

	// Add sorting
	if filter.SortingOrder != nil {
		if strings.ToLower(*filter.SortingOrder) == "asc" {
			query += " ORDER BY m.created_at ASC"
		} else {
			query += " ORDER BY m.created_at DESC"
		}
	} else {
		query += " ORDER BY m.created_at DESC"
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchantItems []model.MerchantItem

	for rows.Next() {
		var merchantItem model.MerchantItem

		var items []byte

		err := rows.Scan(
			&merchantItem.Merchant.ID,
			&merchantItem.Merchant.Name,
			&merchantItem.Merchant.MerchantCategory,
			&merchantItem.Merchant.ImageUrl,
			&merchantItem.Merchant.Location.Lat,
			&merchantItem.Merchant.Location.Long,
			&merchantItem.Merchant.CreatedAt,
			&items,
		)

		if err != nil {
			return nil, err
		}

		if items != nil {
			if err := json.Unmarshal(items, &merchantItem.Items); err != nil {
				log.Err(err)
				return nil, err
			}
		}

		merchantItems = append(merchantItems, merchantItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return merchantItems, nil
}

func (q *Queries) GetMerchantWithItems(ctx context.Context, merchantID int64) (model.MerchantItem, error) {

	var merchantItem model.MerchantItem
	var items []byte

	err := q.db.QueryRowContext(ctx, getMerchantWithItems, merchantID).Scan(
		&merchantItem.Merchant.ID,
		&merchantItem.Merchant.Name,
		&merchantItem.Merchant.MerchantCategory,
		&merchantItem.Merchant.ImageUrl,
		&merchantItem.Merchant.Location.Lat,
		&merchantItem.Merchant.Location.Long,
		&merchantItem.Merchant.CreatedAt,
		&items,
	)

	if err != nil {
		return model.MerchantItem{}, err
	}

	if items != nil {
		if err := json.Unmarshal(items, &merchantItem.Items); err != nil {
			log.Err(err)
			return model.MerchantItem{}, err
		}
	}

	return merchantItem, nil
}

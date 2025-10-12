package repository

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"database/sql"
	"errors"
)

const (
	insertOrderQuery                    = `INSERT INTO orders (user_id, order_estimation_id, is_purchased) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	selectOrderByIdQuery                = `SELECT id, user_id, order_estimation_id, is_purchased, created_at, updated_at FROM orders WHERE id = $1`
	selectOrderByOrderEstimationIdQuery = `SELECT id, user_id, order_estimation_id, is_purchased, created_at, updated_at FROM orders WHERE order_estimation_id = $1`
	selectOrdersByUserAndPurchased      = `SELECT id, user_id, order_estimation_id, is_purchased, created_at, updated_at FROM orders WHERE user_id = $1 AND is_purchased = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	selectUnpurchasedOrderByUserID      = `SELECT id, user_id, order_estimation_id, is_purchased, created_at, updated_at FROM orders WHERE user_id = $1 AND is_purchased = false ORDER BY created_at DESC LIMIT 1`
	updateOrderQuery                    = `UPDATE orders SET order_estimation_id = $2, is_purchased = $3, updated_at = NOW() WHERE id = $1 RETURNING updated_at`

	insertOrderDetailQuery           = `INSERT INTO order_details (order_id, merchant_id, merchant_name, merchant_category, merchant_image_url, merchant_latitude, merchant_longitude) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	selectOrderDetailByIdQuery       = `SELECT id, order_id, merchant_id, merchant_name, merchant_category, merchant_image_url, merchant_latitude, merchant_longitude, created_at, updated_at FROM order_details WHERE id = $1`
	selectOrderDetailByOrderMerchant = `SELECT id, order_id, merchant_id, merchant_name, merchant_category, merchant_image_url, merchant_latitude, merchant_longitude, created_at, updated_at FROM order_details WHERE order_id = $1 AND merchant_id = $2`

	insertOrderItemQuery           = `INSERT INTO order_items (order_detail_id, item_id, item_name, product_category, item_image_url, price_per_item, quantity, total_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	selectOrderItemByDetailAndItem = `SELECT id, order_detail_id, item_id, item_name, product_category, item_image_url, price_per_item, quantity, total_price, created_at, updated_at FROM order_items WHERE order_detail_id = $1 AND item_id = $2`
	updateOrderItemQuery           = `UPDATE order_items SET item_name = $3, product_category = $4, item_image_url = $5, price_per_item = $6, quantity = $7, total_price = $8, updated_at = NOW() WHERE id = $1 AND order_detail_id = $2 RETURNING updated_at`
)

// Order methods
func (q *Queries) CreateOrderWithTx(ctx context.Context, tx *sql.Tx, userID int64, orderEstimationID int64, isPurchased bool) (model.Order, error) {
	db := q.getDB(tx)
	query := insertOrderQuery
	var order model.Order
	err := db.QueryRowContext(ctx, query, userID, orderEstimationID, isPurchased).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, err
	}
	order.UserID = userID
	order.OrderEstimationID = orderEstimationID
	order.IsPurchased = isPurchased
	return order, nil
}

func (q *Queries) CreateOrder(ctx context.Context, userID int64, orderEstimationID int64, isPurchased bool) (model.Order, error) {
	return q.CreateOrderWithTx(ctx, nil, userID, orderEstimationID, isPurchased)
}

func (q *Queries) GetOrderByIDWithTx(ctx context.Context, tx *sql.Tx, id int64) (model.Order, error) {
	db := q.getDB(tx)
	query := selectOrderByIdQuery
	row := db.QueryRowContext(ctx, query, id)
	var order model.Order
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.OrderEstimationID,
		&order.IsPurchased,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, constants.ErrOrderNotFound
		}
		return model.Order{}, err
	}
	return order, nil
}

func (q *Queries) GetOrderByID(ctx context.Context, id int64) (model.Order, error) {
	return q.GetOrderByIDWithTx(ctx, nil, id)
}

func (q *Queries) GetOrderByEstimationIDWithTx(ctx context.Context, tx *sql.Tx, estimationId int64) (model.Order, error) {
	db := q.getDB(tx)
	query := selectOrderByOrderEstimationIdQuery
	row := db.QueryRowContext(ctx, query, estimationId)
	var order model.Order
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.OrderEstimationID,
		&order.IsPurchased,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, constants.ErrOrderNotFound
		}
		return model.Order{}, err
	}
	return order, nil
}

func (q *Queries) GetOrdersByUserIDAndPurchasedWithTx(ctx context.Context, tx *sql.Tx, userID int64, isPurchased bool, limit, offset int) ([]model.Order, error) {
	db := q.getDB(tx)
	query := selectOrdersByUserAndPurchased
	rows, err := db.QueryContext(ctx, query, userID, isPurchased, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderEstimationID,
			&order.IsPurchased,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (q *Queries) GetUnpurchasedOrderByUserIDWithTx(ctx context.Context, tx *sql.Tx, userID int64) (model.Order, error) {
	db := q.getDB(tx)
	query := selectUnpurchasedOrderByUserID
	row := db.QueryRowContext(ctx, query, userID)
	var order model.Order
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.OrderEstimationID,
		&order.IsPurchased,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, constants.ErrNoUnpurchasedOrder
		}
		return model.Order{}, err
	}
	return order, nil
}

func (q *Queries) UpdateOrderWithTx(ctx context.Context, tx *sql.Tx, id int64, orderEstimationID int64, isPurchased bool) (model.Order, error) {
	db := q.getDB(tx)
	query := updateOrderQuery
	var updatedAt sql.NullTime
	err := db.QueryRowContext(ctx, query, id, orderEstimationID, isPurchased).Scan(&updatedAt)
	if err != nil {
		return model.Order{}, err
	}

	// Get the updated order
	return q.GetOrderByIDWithTx(ctx, tx, id)
}

// Order Detail methods
func (q *Queries) CreateOrderDetailWithTx(ctx context.Context, tx *sql.Tx, orderID, merchantID int64, merchantName, merchantCategory, merchantImageURL string, merchantLatitude, merchantLongitude float64) (model.OrderDetail, error) {
	db := q.getDB(tx)
	query := insertOrderDetailQuery
	var orderDetail model.OrderDetail
	err := db.QueryRowContext(ctx, query, orderID, merchantID, merchantName, merchantCategory, merchantImageURL, merchantLatitude, merchantLongitude).Scan(
		&orderDetail.ID,
		&orderDetail.CreatedAt,
		&orderDetail.UpdatedAt,
	)
	if err != nil {
		return model.OrderDetail{}, err
	}
	orderDetail.OrderID = orderID
	orderDetail.MerchantID = merchantID
	orderDetail.MerchantName = merchantName
	orderDetail.MerchantCategory = merchantCategory
	orderDetail.MerchantImageURL = merchantImageURL
	orderDetail.MerchantLatitude = merchantLatitude
	orderDetail.MerchantLongitude = merchantLongitude
	return orderDetail, nil
}

func (q *Queries) GetOrderDetailByIDWithTx(ctx context.Context, tx *sql.Tx, id int64) (model.OrderDetail, error) {
	db := q.getDB(tx)
	query := selectOrderDetailByIdQuery
	row := db.QueryRowContext(ctx, query, id)
	var orderDetail model.OrderDetail
	if err := row.Scan(
		&orderDetail.ID,
		&orderDetail.OrderID,
		&orderDetail.MerchantID,
		&orderDetail.MerchantName,
		&orderDetail.MerchantCategory,
		&orderDetail.MerchantImageURL,
		&orderDetail.MerchantLatitude,
		&orderDetail.MerchantLongitude,
		&orderDetail.CreatedAt,
		&orderDetail.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OrderDetail{}, constants.ErrOrderDetailNotFound
		}
		return model.OrderDetail{}, err
	}
	return orderDetail, nil
}

func (q *Queries) GetOrderDetailByOrderIDAndMerchantIDWithTx(ctx context.Context, tx *sql.Tx, orderID, merchantID int64) (model.OrderDetail, error) {
	db := q.getDB(tx)
	query := selectOrderDetailByOrderMerchant
	row := db.QueryRowContext(ctx, query, orderID, merchantID)
	var orderDetail model.OrderDetail
	if err := row.Scan(
		&orderDetail.ID,
		&orderDetail.OrderID,
		&orderDetail.MerchantID,
		&orderDetail.MerchantName,
		&orderDetail.MerchantCategory,
		&orderDetail.MerchantImageURL,
		&orderDetail.MerchantLatitude,
		&orderDetail.MerchantLongitude,
		&orderDetail.CreatedAt,
		&orderDetail.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OrderDetail{}, constants.ErrOrderDetailNotFound
		}
		return model.OrderDetail{}, err
	}
	return orderDetail, nil
}

// Order Item methods
func (q *Queries) CreateOrderItemWithTx(ctx context.Context, tx *sql.Tx, orderDetailID, itemID int64, itemName, productCategory, itemImageURL string, pricePerItem float64, quantity int32, totalPrice float64) (model.OrderItem, error) {
	db := q.getDB(tx)
	query := insertOrderItemQuery
	var orderItem model.OrderItem
	err := db.QueryRowContext(ctx, query, orderDetailID, itemID, itemName, productCategory, itemImageURL, pricePerItem, quantity, totalPrice).Scan(
		&orderItem.ID,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
	)
	if err != nil {
		return model.OrderItem{}, err
	}
	orderItem.OrderDetailID = orderDetailID
	orderItem.ItemID = itemID
	orderItem.ItemName = itemName
	orderItem.ProductCategory = productCategory
	orderItem.ItemImageURL = itemImageURL
	orderItem.PricePerItem = pricePerItem
	orderItem.Quantity = quantity
	orderItem.TotalPrice = totalPrice
	return orderItem, nil
}

func (q *Queries) GetOrderItemByOrderDetailIDAndItemIDWithTx(ctx context.Context, tx *sql.Tx, orderDetailID, itemID int64) (model.OrderItem, error) {
	db := q.getDB(tx)
	query := selectOrderItemByDetailAndItem
	row := db.QueryRowContext(ctx, query, orderDetailID, itemID)
	var orderItem model.OrderItem
	if err := row.Scan(
		&orderItem.ID,
		&orderItem.OrderDetailID,
		&orderItem.ItemID,
		&orderItem.ItemName,
		&orderItem.ProductCategory,
		&orderItem.ItemImageURL,
		&orderItem.PricePerItem,
		&orderItem.Quantity,
		&orderItem.TotalPrice,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OrderItem{}, constants.ErrOrderItemNotFound
		}
		return model.OrderItem{}, err
	}
	return orderItem, nil
}

func (q *Queries) UpdateOrderItemWithTx(ctx context.Context, tx *sql.Tx, id, orderDetailID int64, itemName, productCategory, itemImageURL string, pricePerItem float64, quantity int32, totalPrice float64) (model.OrderItem, error) {
	db := q.getDB(tx)
	query := updateOrderItemQuery
	var updatedAt sql.NullTime
	err := db.QueryRowContext(ctx, query, id, orderDetailID, itemName, productCategory, itemImageURL, pricePerItem, quantity, totalPrice).Scan(&updatedAt)
	if err != nil {
		return model.OrderItem{}, err
	}

	var orderItem model.OrderItem
	orderItem.ID = id
	orderItem.OrderDetailID = orderDetailID
	orderItem.ItemName = itemName
	orderItem.ProductCategory = productCategory
	orderItem.ItemImageURL = itemImageURL
	orderItem.PricePerItem = pricePerItem
	orderItem.Quantity = quantity
	orderItem.TotalPrice = totalPrice
	orderItem.UpdatedAt = updatedAt.Time
	return orderItem, nil
}

// Overload methods without tx parameter

func (q *Queries) GetOrderByEstimationID(ctx context.Context, estimationId int64) (model.Order, error) {
	return q.GetOrderByEstimationIDWithTx(ctx, nil, estimationId)
}

func (q *Queries) GetOrdersByUserIDAndPurchased(ctx context.Context, userID int64, isPurchased bool, limit, offset int) ([]model.Order, error) {
	return q.GetOrdersByUserIDAndPurchasedWithTx(ctx, nil, userID, isPurchased, limit, offset)
}

func (q *Queries) GetUnpurchasedOrderByUserID(ctx context.Context, userID int64) (model.Order, error) {
	return q.GetUnpurchasedOrderByUserIDWithTx(ctx, nil, userID)
}

func (q *Queries) UpdateOrder(ctx context.Context, id int64, orderEstimationID int64, isPurchased bool) (model.Order, error) {
	return q.UpdateOrderWithTx(ctx, nil, id, orderEstimationID, isPurchased)
}

func (q *Queries) CreateOrderDetail(ctx context.Context, orderID, merchantID int64, merchantName, merchantCategory, merchantImageURL string, merchantLatitude, merchantLongitude float64) (model.OrderDetail, error) {
	return q.CreateOrderDetailWithTx(ctx, nil, orderID, merchantID, merchantName, merchantCategory, merchantImageURL, merchantLatitude, merchantLongitude)
}

func (q *Queries) GetOrderDetailByID(ctx context.Context, id int64) (model.OrderDetail, error) {
	return q.GetOrderDetailByIDWithTx(ctx, nil, id)
}

func (q *Queries) GetOrderDetailByOrderIDAndMerchantID(ctx context.Context, orderID, merchantID int64) (model.OrderDetail, error) {
	return q.GetOrderDetailByOrderIDAndMerchantIDWithTx(ctx, nil, orderID, merchantID)
}

func (q *Queries) CreateOrderItem(ctx context.Context, orderDetailID, itemID int64, itemName, productCategory, itemImageURL string, pricePerItem float64, quantity int32, totalPrice float64) (model.OrderItem, error) {
	return q.CreateOrderItemWithTx(ctx, nil, orderDetailID, itemID, itemName, productCategory, itemImageURL, pricePerItem, quantity, totalPrice)
}

func (q *Queries) GetOrderItemByOrderDetailIDAndItemID(ctx context.Context, orderDetailID, itemID int64) (model.OrderItem, error) {
	return q.GetOrderItemByOrderDetailIDAndItemIDWithTx(ctx, nil, orderDetailID, itemID)
}

func (q *Queries) UpdateOrderItem(ctx context.Context, id, orderDetailID int64, itemName, productCategory, itemImageURL string, pricePerItem float64, quantity int32, totalPrice float64) (model.OrderItem, error) {
	return q.UpdateOrderItemWithTx(ctx, nil, id, orderDetailID, itemName, productCategory, itemImageURL, pricePerItem, quantity, totalPrice)
}

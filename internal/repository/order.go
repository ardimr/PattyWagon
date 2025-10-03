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
func (q *Queries) CreateOrder(ctx context.Context, userID int64, orderEstimationID int64, isPurchased bool) (model.Order, error) {
	query := insertOrderQuery
	var order model.Order
	err := q.db.QueryRowContext(ctx, query, userID, orderEstimationID, isPurchased).Scan(
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

func (q *Queries) GetOrderByID(ctx context.Context, id int64) (model.Order, error) {
	query := selectOrderByIdQuery
	row := q.db.QueryRowContext(ctx, query, id)
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

func (q *Queries) GetOrderByEstimationID(ctx context.Context, estimationId int64) (model.Order, error) {
	query := selectOrderByOrderEstimationIdQuery
	row := q.db.QueryRowContext(ctx, query, estimationId)
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

func (q *Queries) GetOrdersByUserIDAndPurchased(ctx context.Context, userID int64, isPurchased bool, limit, offset int) ([]model.Order, error) {
	query := selectOrdersByUserAndPurchased
	rows, err := q.db.QueryContext(ctx, query, userID, isPurchased, limit, offset)
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

func (q *Queries) GetUnpurchasedOrderByUserID(ctx context.Context, userID int64) (model.Order, error) {
	query := selectUnpurchasedOrderByUserID
	row := q.db.QueryRowContext(ctx, query, userID)
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

func (q *Queries) UpdateOrder(ctx context.Context, id int64, orderEstimationID int64, isPurchased bool) (model.Order, error) {
	query := updateOrderQuery
	var updatedAt sql.NullTime
	err := q.db.QueryRowContext(ctx, query, id, orderEstimationID, isPurchased).Scan(&updatedAt)
	if err != nil {
		return model.Order{}, err
	}

	// Get the updated order
	return q.GetOrderByID(ctx, id)
}

// Order Detail methods
func (q *Queries) CreateOrderDetail(ctx context.Context, orderID, merchantID int64, merchantName, merchantCategory, merchantImageURL string, merchantLatitude, merchantLongitude float64) (model.OrderDetail, error) {
	query := insertOrderDetailQuery
	var orderDetail model.OrderDetail
	err := q.db.QueryRowContext(ctx, query, orderID, merchantID, merchantName, merchantCategory, merchantImageURL, merchantLatitude, merchantLongitude).Scan(
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

func (q *Queries) GetOrderDetailByID(ctx context.Context, id int64) (model.OrderDetail, error) {
	query := selectOrderDetailByIdQuery
	row := q.db.QueryRowContext(ctx, query, id)
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

func (q *Queries) GetOrderDetailByOrderIDAndMerchantID(ctx context.Context, orderID, merchantID int64) (model.OrderDetail, error) {
	query := selectOrderDetailByOrderMerchant
	row := q.db.QueryRowContext(ctx, query, orderID, merchantID)
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
func (q *Queries) CreateOrderItem(ctx context.Context, orderDetailID, itemID int64, itemName, productCategory, itemImageURL string, pricePerItem int64, quantity int32, totalPrice int64) (model.OrderItem, error) {
	query := insertOrderItemQuery
	var orderItem model.OrderItem
	err := q.db.QueryRowContext(ctx, query, orderDetailID, itemID, itemName, productCategory, itemImageURL, pricePerItem, quantity, totalPrice).Scan(
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

func (q *Queries) GetOrderItemByOrderDetailIDAndItemID(ctx context.Context, orderDetailID, itemID int64) (model.OrderItem, error) {
	query := selectOrderItemByDetailAndItem
	row := q.db.QueryRowContext(ctx, query, orderDetailID, itemID)
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

func (q *Queries) UpdateOrderItem(ctx context.Context, id, orderDetailID int64, itemName, productCategory, itemImageURL string, pricePerItem int64, quantity int32, totalPrice int64) (model.OrderItem, error) {
	query := updateOrderItemQuery
	var updatedAt sql.NullTime
	err := q.db.QueryRowContext(ctx, query, id, orderDetailID, itemName, productCategory, itemImageURL, pricePerItem, quantity, totalPrice).Scan(&updatedAt)
	if err != nil {
		return model.OrderItem{}, err
	}

	// Return updated order item (we need to get it since we only have the updated timestamp)
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

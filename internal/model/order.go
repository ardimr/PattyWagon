package model

import (
	"database/sql"
	"time"
)

type Order struct {
	ID                int64     `json:"id" db:"id"`
	UserID            int64     `json:"user_id" db:"user_id"`
	OrderEstimationID int64     `json:"order_estimation_id" db:"order_estimation_id"`
	IsPurchased       bool      `json:"is_purchased" db:"is_purchased"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type OrderDetail struct {
	ID                int64     `json:"id" db:"id"`
	OrderID           int64     `json:"order_id" db:"order_id"`
	MerchantID        int64     `json:"merchant_id" db:"merchant_id"`
	MerchantName      string    `json:"merchant_name" db:"merchant_name"`
	MerchantCategory  string    `json:"merchant_category" db:"merchant_category"`
	MerchantImageURL  string    `json:"merchant_image_url" db:"merchant_image_url"`
	MerchantLatitude  float64   `json:"merchant_latitude" db:"merchant_latitude"`
	MerchantLongitude float64   `json:"merchant_longitude" db:"merchant_longitude"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
	ID              int64     `json:"id" db:"id"`
	OrderDetailID   int64     `json:"order_detail_id" db:"order_detail_id"`
	ItemID          int64     `json:"item_id" db:"item_id"`
	ItemName        string    `json:"item_name" db:"item_name"`
	ProductCategory string    `json:"product_category" db:"product_category"`
	ItemImageURL    string    `json:"item_image_url" db:"item_image_url"`
	PricePerItem    int64     `json:"price_per_item" db:"price_per_item"`
	Quantity        int32     `json:"quantity" db:"quantity"`
	TotalPrice      int64     `json:"total_price" db:"total_price"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type OrderEstimation struct {
	ID                    int64          `json:"id" db:"id"`
	OrderID               sql.NullInt64  `json:"order_id" db:"order_id"`
	UserID                int64          `json:"user_id" db:"user_id"`
	UserLatitude          float64        `json:"user_latitude" db:"user_latitude"`
	UserLongitude         float64        `json:"user_longitude" db:"user_longitude"`
	StartLatitude         float64        `json:"start_latitude" db:"start_latitude"`
	StartLongitude        float64        `json:"start_longitude" db:"start_longitude"`
	TotalDistance         float64        `json:"total_distance" db:"total_distance"`
	EstimatedDeliveryTime float64        `json:"estimated_delivery_time" db:"estimated_delivery_time"`
	OptimalRoutePath      sql.NullString `json:"optimal_route_path" db:"optimal_route_path"`
	TotalItemPrice        int64          `json:"total_item_price" db:"total_item_price"`
	DeliveryFee           int64          `json:"delivery_fee" db:"delivery_fee"`
	TotalPrice            int64          `json:"total_price" db:"total_price"`
	Status                string         `json:"status" db:"status"`
	CreatedAt             time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at" db:"updated_at"`
}

type EstimationRequest struct {
	UserID         int64                   `json:"user_id" validate:"required"`
	UserLatitude   float64                 `json:"user_latitude" validate:"required"`
	UserLongitude  float64                 `json:"user_longitude" validate:"required"`
	StartLatitude  float64                 `json:"start_latitude" validate:"required"`
	StartLongitude float64                 `json:"start_longitude" validate:"required"`
	Items          []EstimationRequestItem `json:"items" validate:"required,min=1"`
}

type EstimationRequestItem struct {
	MerchantID int64 `json:"merchant_id" validate:"required"`
	ItemID     int64 `json:"item_id" validate:"required"`
	Quantity   int32 `json:"quantity" validate:"required,min=1"`
}

type EstimationResponse struct {
	EstimationID   int64 `json:"calculatedEstimateId"`
	TotalItemPrice int64 `json:"estimatedDeliveryTimeInMinutes"`
	TotalPrice     int64 `json:"totalPrice"`
}

type MerchantItemData struct {
	Merchant Merchant
	Item     Item
	Quantity int32
}

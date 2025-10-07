package model

import "time"

type Item struct {
	ID         int64     `db:"id"`
	MerchantID int64     `db:"merchant_id"`
	Name       string    `db:"name"`
	Category   string    `db:"category"`
	Price      float64   `db:"price"`
	ImageURL   string    `db:"image_url"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type FilterItem struct {
	ItemID          int64
	Limit           int
	Offset          int
	Name            string
	ProductCategory string
	CreatedAt       string
}

package model

import "time"

type Item struct {
	ID         int64     `db:"id" json:"id"`
	MerchantID int64     `db:"merchant_id" json:"merchantId"`
	Name       string    `db:"name" json:"name"`
	Category   string    `db:"category" json:"category"`
	Price      float64   `db:"price" json:"price"`
	ImageURL   string    `db:"image_url" json:"imageUrl"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}

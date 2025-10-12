package model

import "time"

type Merchant struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Name      string    `db:"name"`
	Category  string    `db:"category"`
	ImageURL  string    `db:"image_url"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FilterMerchant struct {
	MerchantID       int64
	Limit            int
	Offset           int
	Name             string
	MerchantCategory string
	CreatedAt        string
}

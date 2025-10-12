package model

import "time"

type Cell struct {
	CellID     int64
	Resolution int
}

type Location struct {
	Lat  float64
	Long float64
}

type MerchantLocation struct {
	ID         int64     `db:"id"`
	MerchantID int64     `db:"merchant_id"`
	H3Index    int64     `db:"h3_index"`
	Resolution int       `db:"resolution"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

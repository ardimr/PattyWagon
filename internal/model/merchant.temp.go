package model

import "time"

type Merchant struct {
	ID               int64
	Location         Location
	Name             string
	MerchantCategory string
	ImageUrl         string
	CreatedAt        time.Time
}

type ListMerchantWithItemParams struct {
	Cell *Cell
	MerchantParams
}

type MerchantParams struct {
	MerchantID       *int64
	Limit            int
	Offset           int
	Name             *string
	MerchantCategory *string
}

type MerchantLocation struct {
	ID         int64
	MerchantID int64
	H3Index    int64
	Resolution int8
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MerchantCategory string

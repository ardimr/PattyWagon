package model

import "time"

type OrderEstimation struct {
	UserLocation Location
	Orders       []Order
}

type Order struct {
	MerchantID      int64
	IsStartingPoint bool
	Items           []OrderItem
}

type Item struct {
	ID              int64
	Name            string
	Price           int64
	ProductCategory string
	ImageUrl        string
	CreatedAt       time.Time
}

type OrderItem struct {
	ItemID   int64
	Quantity int
}

type EstimationPrice struct {
	ID                         int64
	EstimatedDeliveryInMinutes int64
	TotalPrice                 int64
}

type FindNerbyMerchantParams struct {
	UserLocation Location
	MerchantParams
}

type MerchantItem struct {
	Merchant Merchant
	Items    []Item
}

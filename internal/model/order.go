package model

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
	ID    int64
	Price int64
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
	UserLocation     Location
	MerchantID       *int64
	Limit            *int
	Offset           *int
	Name             *string
	MerchantCategory *string
}

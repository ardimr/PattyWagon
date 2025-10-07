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

// type Item struct {
// 	ID              int64     `json:"id"`
// 	Name            string    `json:"name"`
// 	Price           float64   `json:"price"`
// 	ProductCategory string    `json:"category"`
// 	ImageUrl        string    `json:"image_url"`
// 	CreatedAt       time.Time `json:"created_at"`
// }

type OrderItem struct {
	ItemID   int64
	Quantity int
}

type EstimationPrice struct {
	ID                         int64
	EstimatedDeliveryInMinutes int64
	TotalPrice                 float64
}

type FindNerbyMerchantParams struct {
	UserLocation Location
	MerchantParams
}

type MerchantItem struct {
	Merchant Merchant
	Items    []Item
}

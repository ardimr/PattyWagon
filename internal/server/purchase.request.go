package server

import (
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
)

type OrderEstimationRequest struct {
	UserLocation LocationRequest `json:"userLocation" validate:"required"`
	Orders       []OrderRequest  `json:"orders"`
}

type LocationRequest struct {
	Lat  float64 `json:"lat" validate:"required,latitude"`
	Long float64 `json:"long" validate:"required,longitude"`
}

type OrderRequest struct {
	MerchantID      string             `json:"merchantId" validate:"required"`
	IsStartingPoint bool               `json:"isStartingPoint" validate:"required"`
	Items           []OrderItemRequest `json:"items" validate:"required"`
}

type OrderItemRequest struct {
	ItemID   string `json:"itemId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}

func (r *OrderItemRequest) ToModel() model.OrderItem {
	return model.OrderItem{
		ItemID:   utils.ConvertIDString2Int(r.ItemID, 0),
		Quantity: r.Quantity,
	}
}

func (r *LocationRequest) ToModel() model.Location {
	return model.Location{
		Lat:  r.Lat,
		Long: r.Long,
	}
}

func (r *OrderRequest) ToModel() model.Order {
	merchantId := utils.ConvertIDString2Int(r.MerchantID, 0)

	var orderItems []model.OrderItem
	for _, item := range r.Items {
		orderItems = append(orderItems, item.ToModel())
	}

	return model.Order{
		MerchantID:      merchantId,
		IsStartingPoint: r.IsStartingPoint,
		Items:           orderItems,
	}
}

func (r *OrderEstimationRequest) ToModel() model.OrderEstimation {
	var res model.OrderEstimation
	var orders []model.Order

	for _, order := range r.Orders {
		orders = append(orders, order.ToModel())
	}

	res.UserLocation = r.UserLocation.ToModel()
	res.Orders = orders
	return res
}

package server

import (
	"PattyWagon/internal/model"
	"strconv"
)

type EstimationPriceResponse struct {
	CalculateEstimateID        string `json:"calculateEstimateId"`
	TotalPrice                 int64  `json:"totalPrice"`
	EstimatedDeliveryInMinutes int64  `json:"estimatedDeliveryInMinutes"`
}

func (r *EstimationPriceResponse) FromModel(input model.EstimationPrice) {
	r.CalculateEstimateID = strconv.Itoa(int(input.ID))
	r.EstimatedDeliveryInMinutes = input.EstimatedDeliveryInMinutes
	r.TotalPrice = input.TotalPrice
}

func NewEstimationPriceResponse(input model.EstimationPrice) EstimationPriceResponse {
	return EstimationPriceResponse{
		CalculateEstimateID:        strconv.Itoa(int(input.ID)),
		TotalPrice:                 input.TotalPrice,
		EstimatedDeliveryInMinutes: input.EstimatedDeliveryInMinutes,
	}
}

type LocationResponse struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Merchant struct {
	MerchantID       string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory string           `json:"merchantCategory"`
	ImageUrl         string           `json:"imageUrl"`
	Location         LocationResponse `json:"location"`
	CreatedAt        int64            `json:"createdAt"`
}

type Item struct {
	ItemID          string `json:"itemId"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	ImageUrl        string `json:"imageUrl"`
	CreatedAt       int64  `json:"createdAt"`
}

type MerchantWithItem struct {
	Merchant Merchant `json:"merchant"`
	Items    []Item   `json:"item"`
}

type FindNearbyMerchantsResponse struct {
	Data []MerchantWithItem `json:"data"`
}

func NewLocationResponse(input model.Location) LocationResponse {
	return LocationResponse{
		Lat:  input.Lat,
		Long: input.Long,
	}
}

func NewMerchantResponse(input model.Merchant) Merchant {
	return Merchant{
		MerchantID:       strconv.Itoa(int(input.ID)),
		Location:         NewLocationResponse(input.Location),
		Name:             input.Name,
		MerchantCategory: input.MerchantCategory,
		ImageUrl:         input.ImageUrl,
		CreatedAt:        input.CreatedAt.Unix(),
	}
}

func NewItemResponse(input model.Item) Item {
	return Item{
		ItemID:          strconv.Itoa(int(input.ID)),
		Name:            input.Name,
		ProductCategory: input.ProductCategory,
		ImageUrl:        input.ImageUrl,
		CreatedAt:       input.CreatedAt.Unix(),
	}
}

func NewMultipleItemsResponse(inputs []model.Item) []Item {
	var items []Item
	for _, item := range inputs {
		items = append(items, NewItemResponse(item))
	}
	return items
}

func NewMerchantWithItem(input model.MerchantItem) MerchantWithItem {
	return MerchantWithItem{
		Merchant: NewMerchantResponse(input.Merchant),
		Items:    NewMultipleItemsResponse(input.Items),
	}
}

func NewFindNearbyMerchantsResponse(inputs []model.MerchantItem) FindNearbyMerchantsResponse {
	var merchantWithItems []MerchantWithItem

	for _, input := range inputs {
		merchantWithItems = append(merchantWithItems, NewMerchantWithItem(input))
	}
	return FindNearbyMerchantsResponse{
		Data: merchantWithItems,
	}

}

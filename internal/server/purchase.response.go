package server

import (
	"PattyWagon/internal/model"
	"strconv"
	"time"
)

type EstimationPriceResponse struct {
	CalculateEstimateID        string  `json:"calculateEstimateId"`
	TotalPrice                 float64 `json:"totalPrice"`
	EstimatedDeliveryInMinutes int64   `json:"estimatedDeliveryInMinutes"`
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
	CreatedAt        time.Time        `json:"createdAt"`
}

type Item struct {
	ItemID          string    `json:"itemId"`
	Name            string    `json:"name"`
	ProductCategory string    `json:"productCategory"`
	ImageUrl        string    `json:"imageUrl"`
	CreatedAt       time.Time `json:"createdAt"`
}

type MerchantWithItem struct {
	Merchant Merchant `json:"merchant"`
	Items    []Item   `json:"item"`
}

type FindNearbyMerchantsResponseMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
type FindNearbyMerchantsResponse struct {
	Data []MerchantWithItem              `json:"data"`
	Meta FindNearbyMerchantsResponseMeta `json:"meta"`
}

func NewLocationResponse(lat, long float64) LocationResponse {
	return LocationResponse{
		Lat:  lat,
		Long: long,
	}
}

func NewMerchantResponse(input model.Merchant) Merchant {
	return Merchant{
		MerchantID:       strconv.Itoa(int(input.ID)),
		Location:         NewLocationResponse(input.Latitude, input.Longitude),
		Name:             input.Name,
		MerchantCategory: *input.Category,
		ImageUrl:         input.ImageURL,
		CreatedAt:        input.CreatedAt,
	}
}

func NewItemResponse(input model.Item) Item {
	return Item{
		ItemID:          strconv.Itoa(int(input.ID)),
		Name:            input.Name,
		ProductCategory: input.Category,
		ImageUrl:        input.ImageURL,
		CreatedAt:       input.CreatedAt,
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

func NewFindNearbyMerchantsResponse(inputs []model.MerchantItem, meta FindNearbyMerchantsResponseMeta) FindNearbyMerchantsResponse {
	var merchantWithItems []MerchantWithItem

	for _, input := range inputs {
		merchantWithItems = append(merchantWithItems, NewMerchantWithItem(input))
	}
	return FindNearbyMerchantsResponse{
		Data: merchantWithItems,
		Meta: meta,
	}

}

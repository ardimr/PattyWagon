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

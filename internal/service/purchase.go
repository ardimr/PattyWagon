package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"PattyWagon/logger"
	"context"

	"github.com/kwahome/go-haversine/pkg/haversine"
)

func (s *Service) FindNearbyMerchants(ctx context.Context, searchParams model.FindNerbyMerchantParams) ([]model.Merchant, error) {
	log := logger.GetLoggerFromContext(ctx)

	var merchants []model.Merchant

	neighborCells, err := s.locationService.FindNearby(ctx, searchParams.UserLocation)
	if err != nil {
		return []model.Merchant{}, err
	}

	for _, cell := range neighborCells {
		merchant, err := s.repository.GetMerchantByCellID(ctx, cell.CellID)
		if err != nil {
			return nil, err
		}

		merchants = append(merchants, merchant)
	}

	log.Printf("Total merchants: %d", len(merchants))
	return merchants, nil
}

func (s *Service) EstimateOrderPrice(ctx context.Context, orderEstimation model.OrderEstimation) (model.EstimationPrice, error) {
	isStartingPointValid := false

	var locations []model.Location
	var totalPrice int64 = 0
	for _, order := range orderEstimation.Orders {
		isStartingPointValid = isStartingPointValid != order.IsStartingPoint

		merchant, err := s.repository.GetMerchantByID(ctx, order.MerchantID)
		if err != nil {
			return model.EstimationPrice{}, err
		}

		if !s.validateDistance(ctx, orderEstimation.UserLocation, merchant.Location) {
			return model.EstimationPrice{}, constants.ErrMerchantTooFar
		}

		for _, orderItem := range order.Items {
			item, err := s.repository.GetItemByID(ctx, orderItem.ItemID)
			if err != nil {
				return model.EstimationPrice{}, err
			}

			totalPrice += int64(orderItem.Quantity) * item.Price
		}

		merchant.Location.IsStartingPoint = order.IsStartingPoint
		locations = append(locations, merchant.Location)
	}

	if !isStartingPointValid {
		return model.EstimationPrice{}, constants.ErrInvalidStartingPoint
	}

	estimatedDeliveryTimeInMinutes, err := s.locationService.EstimateDeliveryTimeInMinutes(ctx, locations)
	if err != nil {
		return model.EstimationPrice{}, err
	}

	estimationPrice := model.EstimationPrice{
		EstimatedDeliveryInMinutes: estimatedDeliveryTimeInMinutes,
		TotalPrice:                 totalPrice,
	}

	return estimationPrice, nil
}

func (s *Service) validateDistance(ctx context.Context, user, merchant model.Location) bool {
	userCoord := haversine.Coordinate{
		Latitude:  user.Lat,
		Longitude: user.Long,
	}

	merchantCoord := haversine.Coordinate{
		Latitude:  merchant.Lat,
		Longitude: merchant.Long,
	}

	distance := merchantCoord.DistanceTo(userCoord, haversine.M)

	return distance < 3000
}

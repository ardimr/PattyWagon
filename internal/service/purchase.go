package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"

	"github.com/kwahome/go-haversine/pkg/haversine"
)

func (s *Service) EstimateOrderPrice(ctx context.Context, orderEstimation model.OrderEstimation) (model.EstimationPrice, error) {
	isStartingPointValid := false

	for _, order := range orderEstimation.Orders {
		isStartingPointValid = isStartingPointValid != order.IsStartingPoint

		merchant, err := s.repository.GetMerchantByID(ctx, order.MerchantID)
		if err != nil {
			return model.EstimationPrice{}, err
		}

		if !s.validateDistance(ctx, orderEstimation.UserLocation, merchant.Location) {
			return model.EstimationPrice{}, constants.ErrMerchantTooFar
		}
	}

	if !isStartingPointValid {
		return model.EstimationPrice{}, constants.ErrInvalidStartingPoint
	}

	return model.EstimationPrice{}, nil
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

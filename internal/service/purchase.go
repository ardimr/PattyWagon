package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"PattyWagon/logger"
	"context"

	"github.com/kwahome/go-haversine/pkg/haversine"
)

func (s *Service) FindNearbyMerchants(ctx context.Context, userLocation model.Location, searchParams model.FindNerbyMerchantParams) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)

	merchants, err := s.findNearbyMerchantsWithStrategy(ctx, userLocation, searchParams)
	if err != nil {
		return nil, err
	}

	log.Printf("total merchants: %d", len(merchants))
	return merchants, nil
}

func (s *Service) findNearbyMerchantsWithStrategy(ctx context.Context, userLocation model.Location, searchParams model.FindNerbyMerchantParams) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)

	var merchants []model.MerchantItem
	numAcquiredMerchants := 0
	searchingLevel := 1
	numRequiredMerchants := searchParams.MerchantParams.Limit + searchParams.MerchantParams.Offset

	cellMap := make(map[int64]model.Cell, 0)

	for numAcquiredMerchants < numRequiredMerchants {
		log.Printf("Finding nearby merchants: (%d/%d)", numAcquiredMerchants, numRequiredMerchants)

		unseenCells := make([]model.Cell, 0)

		neighborCells, err := s.locationService.FindNearby(ctx, userLocation, searchingLevel)
		if err != nil {
			return nil, err
		}

		log.Printf("neighbor cells: %d", len(neighborCells))
		for _, cell := range neighborCells {
			if _, exists := cellMap[cell.ID]; !exists {
				unseenCells = append(unseenCells, cell)
			}
		}

		queryParams := model.ListMerchantWithItemParams{
			Cells:          unseenCells,
			MerchantParams: searchParams.MerchantParams,
		}

		unseenMerchants, err := s.repository.ListMerchantWithItems(ctx, queryParams)
		if err != nil {
			return nil, err
		}

		numAcquiredMerchants += len(unseenMerchants)
		searchingLevel += 1
		merchants = append(merchants, unseenMerchants...)
	}

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

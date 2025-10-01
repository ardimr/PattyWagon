package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
	"PattyWagon/logger"
	"context"
	"os"

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
	numRequiredMerchants := searchParams.MerchantParams.Limit + searchParams.MerchantParams.Offset

	cellMap := make(map[int64]model.Cell, 0)
	seenMerchants := make(map[int64]struct{}, 0)

	resolution := utils.String2Int32(os.Getenv("H3_START_RESOLUTION"), 8)
	k := 1

	// Key Strategy
	// 0. If the total merchants is less than 2 * numRequiredMerchants then just direct query to database
	// 1. find nearby with decreasing h3 resolution
	// 2. if we have reached the lowest h3 resolution (0) and number of required merchants is still not fulfilled, start using k-ring approach

	for (numAcquiredMerchants < numRequiredMerchants) && (resolution >= 0 || k <= 15) {
		log.Printf("Finding nearby merchants: (%d/%d)", numAcquiredMerchants, numRequiredMerchants)
		if resolution >= 0 {
			filteredMerchants, err := s.findNearbyMerchantsByResolution(ctx, userLocation, searchParams, resolution, seenMerchants)
			if err != nil {
				return nil, err
			}
			numAcquiredMerchants += len(filteredMerchants)
			merchants = append(merchants, filteredMerchants...)
			resolution -= 1
		} else {
			filteredMerchants, err := s.findNearbyMerchantByKRing(ctx, userLocation, searchParams, k, seenMerchants, cellMap)
			if err != nil {
				return nil, err
			}
			numAcquiredMerchants += len(filteredMerchants)
			merchants = append(merchants, filteredMerchants...)
			k += 1
		}
	}

	if numAcquiredMerchants > numRequiredMerchants {
		return merchants[searchParams.Offset:numRequiredMerchants], nil
	}
	return merchants, nil
}

func (s *Service) findNearbyMerchantsByResolution(ctx context.Context, userLocation model.Location, searchParams model.FindNerbyMerchantParams, resolution int, seenMerchants map[int64]struct{}) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)

	var merchants []model.MerchantItem

	log.Printf("Resolution: %d ", resolution)
	cell, err := s.locationService.FindCellIDByResolution(ctx, userLocation, resolution)
	if err != nil {
		return nil, err
	}

	// Get all merchants within the cell that satisfy the filter
	filter := model.ListMerchantWithItemParams{
		Cell:           cell,
		MerchantParams: searchParams.MerchantParams,
	}

	merchantItems, err := s.repository.ListMerchantWithItems(ctx, filter)
	if err != nil {
		return nil, err
	}

	for _, merchant := range merchantItems {
		if _, exists := seenMerchants[merchant.Merchant.ID]; !exists {
			seenMerchants[merchant.Merchant.ID] = struct{}{}
			merchants = append(merchants, merchant)
		}
	}
	return merchants, nil
}

func (s *Service) findNearbyMerchantByKRing(ctx context.Context, userLocation model.Location, searchParams model.FindNerbyMerchantParams, k int, seenMerchants map[int64]struct{}, cellMap map[int64]model.Cell) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)
	log.Printf("K-ring: %d ", k)

	var merchants []model.MerchantItem
	unseenCells := make([]model.Cell, 0)

	cells, err := s.locationService.FindKRingCellIDs(ctx, userLocation, k)
	if err != nil {
		return nil, err
	}

	for _, cell := range cells {
		if _, exists := cellMap[cell.CellID]; !exists {
			unseenCells = append(unseenCells, cell)
		}
	}

	for _, cell := range unseenCells {
		cellMap[cell.CellID] = cell
		queryParams := model.ListMerchantWithItemParams{
			Cell:           cell,
			MerchantParams: searchParams.MerchantParams,
		}

		filteredMerchants, err := s.repository.ListMerchantWithItems(ctx, queryParams)
		if err != nil {
			return nil, err
		}

		for _, merchant := range filteredMerchants {
			if _, exists := seenMerchants[merchant.Merchant.ID]; !exists {
				seenMerchants[merchant.Merchant.ID] = struct{}{}
				merchants = append(merchants, merchant)
			}
		}
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

	// estimatedDeliveryTimeInMinutes, err := s.locationService.EstimateDeliveryTimeInMinutes(ctx, locations)
	// if err != nil {
	// 	return model.EstimationPrice{}, err
	// }

	estimationPrice := model.EstimationPrice{
		EstimatedDeliveryInMinutes: 6,
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

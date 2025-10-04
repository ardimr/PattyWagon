package service

import (
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
	"PattyWagon/logger"
	"context"
	"os"
	"slices"
	"sync"
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

func (s *Service) findNearbyMerchantsWithStrategy(ctx context.Context, userLocation model.Location, filter model.FindNerbyMerchantParams) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)

	var merchants []model.MerchantItem

	numAcquiredMerchants := 0
	numRequiredMerchants := filter.MerchantParams.Limit + filter.MerchantParams.Offset
	log.Printf("limit: %d offset:%d requiredMerchants: %d", filter.Limit, filter.Offset, numRequiredMerchants)

	cellMap := make(map[int64]model.Cell, 0)
	seenMerchants := make(map[int64]struct{}, 0)

	resolution := utils.String2Int32(os.Getenv("H3_START_RESOLUTION"), 8)
	kRing := 1
	maxKRing := 30

	// Key Strategy
	// - If the total merchants is less than 2 * numRequiredMerchants then just direct query to database
	// - find nearby merhants starting from the k-ring 1
	// - if the retrieved merchants less than numRequiredMerchants, expand to k-ring
	// - if k-ring reaches limit (max k-ring) just return all merchants from database ordered by distance

	// Precheck

	if filter.MerchantID != nil {
		log.Printf("Get merchant with items directly: (%d)", *filter.MerchantID)
		merchantItem, err := s.repository.GetMerchantWithItems(ctx, *filter.MerchantID)
		if err != nil {
			return nil, err
		}
		return []model.MerchantItem{merchantItem}, nil
	}

	// Phase 1
	for (numAcquiredMerchants < numRequiredMerchants) && (kRing < maxKRing) {
		log.Printf("Finding nearby merchants: (%d/%d)", numAcquiredMerchants, numRequiredMerchants)

		filteredMerchants, err := s.findNearbyMerchantsByKRing(ctx, userLocation, filter.MerchantParams, resolution, kRing, seenMerchants, cellMap)
		if err != nil {
			return nil, err
		}

		numAcquiredMerchants += len(filteredMerchants)
		merchants = append(merchants, filteredMerchants...)
		kRing += 1
	}

	// Phase 2
	if numAcquiredMerchants < numRequiredMerchants {
		log.Println("Total acquired merchants less than threshold -> find from database directly")
		filteredMerchants, err := s.findNearbyMerchantsFromDatabase(ctx, filter.MerchantParams, seenMerchants)
		if err != nil {
			return nil, err
		}

		numAcquiredMerchants += len(filteredMerchants)
		merchants = append(merchants, filteredMerchants...)
	}

	log.Printf("unsorted merchants: %d", len(merchants))
	return s.sortAndLimitNearbyMerchants(userLocation, merchants, filter.Offset, filter.Limit), nil
}

func (s *Service) findNearbyMerchantsByKRing(ctx context.Context, userLocation model.Location, filter model.MerchantParams, resolution, k int, seenMerchants map[int64]struct{}, cellMap map[int64]model.Cell) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)
	log.Printf("K-ring: %d ", k)

	var merchants []model.MerchantItem
	var lock sync.Mutex
	var wg sync.WaitGroup
	unseenCells := make([]model.Cell, 0)

	cells, err := s.locationService.FindKRingCellIDs(ctx, userLocation, resolution, k)
	if err != nil {
		return nil, err
	}

	for _, cell := range cells {
		if _, exists := cellMap[cell.CellID]; !exists {
			unseenCells = append(unseenCells, cell)
		}
	}

	// log.Printf("unseen cells: %v", unseenCells)

	// TODO: implement concurrent query
	errCh := make(chan error, len(unseenCells))

	for _, cell := range unseenCells {
		cellMap[cell.CellID] = cell
		queryParams := model.ListMerchantWithItemParams{
			Cell:           &cell,
			MerchantParams: filter,
		}

		wg.Add(1)

		go func(c model.Cell, params model.ListMerchantWithItemParams) {
			defer wg.Done()
			filteredMerchants, err := s.repository.ListMerchantWithItems(ctx, params)
			if err != nil {
				errCh <- err
			}

			lock.Lock()
			defer lock.Unlock()

			for _, merchant := range filteredMerchants {
				if _, exists := seenMerchants[merchant.Merchant.ID]; !exists {
					seenMerchants[merchant.Merchant.ID] = struct{}{}
					merchants = append(merchants, merchant)
				}
			}
		}(cell, queryParams)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return merchants, nil
}

func (s *Service) findNearbyMerchantsFromDatabase(ctx context.Context, filter model.MerchantParams, seenMerchants map[int64]struct{}) ([]model.MerchantItem, error) {
	log := logger.GetLoggerFromContext(ctx)
	log.Printf("Searching merchants from database with filter: %+v\n", filter)

	var merchants []model.MerchantItem
	queryParams := model.ListMerchantWithItemParams{
		MerchantParams: filter,
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

	return merchants, nil
}

func (s *Service) sortAndLimitNearbyMerchants(userLocation model.Location, merchants []model.MerchantItem, offset, limit int) []model.MerchantItem {
	lenMerchants := len(merchants)

	slices.SortFunc(merchants, func(m1, m2 model.MerchantItem) int {
		d1 := utils.CalculateDistance(userLocation.Lat, userLocation.Long, m1.Merchant.Latitude, m1.Merchant.Longitude)
		d2 := utils.CalculateDistance(userLocation.Lat, userLocation.Long, m2.Merchant.Latitude, m2.Merchant.Longitude)
		return int(d1 - d2)
	})

	if offset >= lenMerchants {
		return merchants
	}

	end := offset + limit
	if end > lenMerchants {
		end = len(merchants)
	}

	return merchants[offset:end]
}

// func (s *Service) EstimateOrderPrice(ctx context.Context, orderEstimation model.OrderEstimation) (model.EstimationPrice, error) {
// 	isStartingPointValid := false

// 	var locations []model.Location
// 	var totalPrice float64 = 0
// 	for _, order := range orderEstimation.Orders {
// 		isStartingPointValid = isStartingPointValid != order.IsStartingPoint

// 		merchant, err := s.repository.GetMerchantByID(ctx, order.MerchantID)
// 		if err != nil {
// 			return model.EstimationPrice{}, err
// 		}

// 		if !s.validateDistance(orderEstimation.UserLocation, merchant.Location) {
// 			return model.EstimationPrice{}, constants.ErrMerchantTooFar
// 		}

// 		for _, orderItem := range order.Items {
// 			item, err := s.repository.GetItemByID(ctx, orderItem.ItemID)
// 			if err != nil {
// 				return model.EstimationPrice{}, err
// 			}

// 			totalPrice += float64(orderItem.Quantity) * item.Price
// 		}

// 		merchant.Location.IsStartingPoint = order.IsStartingPoint
// 		locations = append(locations, merchant.Location)
// 	}

// 	if !isStartingPointValid {
// 		return model.EstimationPrice{}, constants.ErrInvalidStartingPoint
// 	}

// 	// estimatedDeliveryTimeInMinutes, err := s.locationService.EstimateDeliveryTimeInMinutes(ctx, locations)
// 	// if err != nil {
// 	// 	return model.EstimationPrice{}, err
// 	// }

// 	estimationPrice := model.EstimationPrice{
// 		EstimatedDeliveryInMinutes: 6,
// 		TotalPrice:                 totalPrice,
// 	}

// 	return estimationPrice, nil
// }

func (s *Service) validateDistance(user, merchant model.Location) bool {

	distance := utils.CalculateDistance(user.Lat, user.Long, merchant.Lat, merchant.Long)

	return distance < 3000
}

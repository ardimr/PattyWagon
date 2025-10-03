package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/types"
	"context"
	"math"
	"sync"
)

type Location struct {
	Name string
	Lat  float64
	Lon  float64
}

type merchantResult struct {
	location Location
	err      error
	index    int
}

func (s *Service) FindOptimalRoute(ctx context.Context, startLat, startLon, userLat, userLon float64, merchantIDs []int64) (types.RouteResult, error) {
	// Create start and user locations
	startLocation := Location{
		Name: "Start",
		Lat:  startLat,
		Lon:  startLon,
	}

	userLocation := Location{
		Name: "User",
		Lat:  userLat,
		Lon:  userLon,
	}

	merchantLocations, err := s.getMerchantsLocationsConurrently(ctx, merchantIDs)
	if err != nil {
		return types.RouteResult{}, err
	}

	// Solve TSP to find optimal route
	result := solveTSP(startLocation, userLocation, merchantLocations)

	return result, nil
}

func (s *Service) CalculateEstimatedDeliveryTime(ctx context.Context, route types.RouteResult) (float64, error) {
	// The delivery time is already calculated in the RouteResult
	return route.DeliveryTime, nil
}

func (s *Service) getMerchantsLocationsConurrently(ctx context.Context, merchantIDs []int64) ([]Location, error) {
	if len(merchantIDs) == 0 {
		return []Location{}, nil
	}

	resultChan := make(chan merchantResult, len(merchantIDs))

	var wg sync.WaitGroup
	for i, merchantID := range merchantIDs {
		wg.Add(1)
		go func(index int, id int64) {
			defer wg.Done()

			merchant, err := s.GetMerchant(ctx, id)
			if err != nil {
				resultChan <- merchantResult{
					err:   err,
					index: index,
				}
				return
			}

			location := Location{
				Name: merchant.Name,
				Lat:  merchant.Latitude,
				Lon:  merchant.Longitude,
			}

			resultChan <- merchantResult{
				location: location,
				index:    index,
			}
		}(i, merchantID)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	merchantLocations := make([]Location, len(merchantIDs))
	for result := range resultChan {
		if result.err != nil {
			return nil, constants.WrapError(constants.ErrFailedToGetMerchant, result.err)
		}
		merchantLocations[result.index] = result.location
	}

	return merchantLocations, nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	deltaLat := radians(lat2 - lat1)
	deltaLon := radians(lon2 - lon1)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + (math.Cos(radians(lat1)) * math.Cos(radians(lat2)) * math.Sin(deltaLon/2) * math.Sin(deltaLon/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return d
}

func radians(deg float64) float64 {
	return deg * math.Pi / 180
}

func toDistanceGraph(locations []Location) map[string]map[string]float64 {
	n := len(locations)
	graph := make(map[string]map[string]float64)
	for i := 0; i < n; i++ {
		dist := make(map[string]float64)
		for j := 0; j < n; j++ {
			if i != j {
				dist[locations[j].Name] = haversine(locations[i].Lat, locations[i].Lon, locations[j].Lat, locations[j].Lon)
			}
		}
		graph[locations[i].Name] = dist
	}
	return graph
}

func permutations(arr []string) [][]string {
	var helper func([]string, int)
	res := [][]string{}

	helper = func(a []string, i int) {
		if i == len(a)-1 {
			tmp := make([]string, len(a))
			copy(tmp, a)
			res = append(res, tmp)
			return
		}
		for j := i; j < len(a); j++ {
			a[i], a[j] = a[j], a[i]
			helper(a, i+1)
			a[i], a[j] = a[j], a[i]
		}
	}

	helper(arr, 0)
	return res
}

func solveTSP(startLocation, userLocation Location, merchantLocations []Location) types.RouteResult {
	var merchant []string
	locations := []Location{startLocation}
	for _, loc := range merchantLocations {
		locations = append(locations, loc)
		merchant = append(merchant, loc.Name)
	}
	locations = append(locations, userLocation)

	graphDistance := toDistanceGraph(locations)

	start := startLocation.Name
	destination := userLocation.Name

	minCost := math.MaxFloat64
	var bestPath []string

	for _, perm := range permutations(merchant) {
		cost := 0.0
		valid := true

		for i := 0; i < len(perm)-1; i++ {
			from, to := perm[i], perm[i+1]
			if c, ok := graphDistance[from][to]; ok {
				cost += c
			} else {
				valid = false
				break
			}
		}

		// Add start merchant
		first := perm[0]
		if c, ok := graphDistance[start][first]; ok && valid {
			cost += c
		} else {
			valid = false
		}

		// Add final leg to destination
		last := perm[len(perm)-1]
		if c, ok := graphDistance[last][destination]; ok && valid {
			cost += c
		} else {
			valid = false
		}

		if valid && cost < minCost {
			minCost = cost
			bestPath = append([]string{}, start)
			bestPath = append(bestPath, perm...)
			bestPath = append(bestPath, destination)
		}
	}

	deliveryTime := minCost / 2400.0

	return types.RouteResult{
		Path:         bestPath,
		TotalCost:    minCost,
		DeliveryTime: deliveryTime,
	}
}

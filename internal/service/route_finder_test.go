package service

import (
	"context"
	"testing"
)

// Test for the core TSP algorithm with known data
func TestSolveTSP(t *testing.T) {
	startLocation := Location{"Start", 22.1234, 12.5678}
	userLocation := Location{"User", 22.1234, -11.5678}

	merchantLocations := []Location{
		{"Merchant A", 40.7128, -74.0060},
		{"Merchant B", 37.1234, -122.6543},
		{"Merchant C", -12.8756, 45.1234},
		{"Merchant D", 51.5076, -0.1227},
	}

	result := solveTSP(startLocation, userLocation, merchantLocations)

	if len(result.Path) == 0 {
		t.Error("Expected non-empty path")
	}

	if result.TotalCost <= 0 {
		t.Error("Expected positive total cost")
	}

	if result.DeliveryTime <= 0 {
		t.Error("Expected positive delivery time")
	}

	// Should include all merchants plus start and end
	expectedPathLength := len(merchantLocations) + 2
	if len(result.Path) != expectedPathLength {
		t.Errorf("Expected path length %d, got %d", expectedPathLength, len(result.Path))
	}

	t.Logf("Best path: %v", result.Path)
	t.Logf("Total cost: %f", result.TotalCost)
	t.Logf("Delivery time: %f", result.DeliveryTime)
}

// Test concurrent merchant location fetching
func TestGetMerchantsLocationsConcurrently(t *testing.T) {
	// This test would require mocking the merchant service
	// For now, just test the structure
	service := &Service{}
	merchantIDs := []int64{}

	locations, err := service.getMerchantsLocationsConurrently(context.Background(), merchantIDs)

	if err != nil {
		t.Errorf("Expected no error for empty merchant IDs, got: %v", err)
	}

	if len(locations) != 0 {
		t.Errorf("Expected empty locations for empty merchant IDs, got: %d", len(locations))
	}

	t.Log("Concurrent merchant location fetching structure is correct")
}

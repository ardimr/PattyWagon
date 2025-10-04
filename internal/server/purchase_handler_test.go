package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/database"
	imagecompressor "PattyWagon/internal/image_compressor"
	"PattyWagon/internal/location"
	mocklocationservice "PattyWagon/internal/mock_location_service"
	"PattyWagon/internal/mock_repository"
	"PattyWagon/internal/model"
	"PattyWagon/internal/repository"
	"PattyWagon/internal/service"
	"PattyWagon/internal/storage"
	"PattyWagon/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testPurchaseSetup(t *testing.T) *Server {
	// repo := &mock_repository.TestRepositoryMock{}
	repo := repository.New(database.New(
		"localhost",
		"5432",
		"patty-wagon-dev",
		"postgres",
		"postgres",
		"public",
		&database.ConnectionPoolConfig{
			MaxOpenConns:    20,
			MaxIdleConns:    10,
			ConnMaxIdleTime: 30 * time.Second,
			ConnMaxLifeTime: 300 * time.Second,
		},
	))
	storage := storage.New("localhost:9000", "team-solid", "@team-solid", storage.Option{MaxConcurrent: 5})
	imageCompressor := imagecompressor.New(5, 50)
	// locationSvc := &mocklocationservice.MockLocationService{}
	locationSvc := location.NewService()
	svc := service.New(repo, storage, imageCompressor, locationSvc)

	// testPopulateMockRepo(t, repo)
	// testPopulateMockLocationService(t, locationSvc)
	return &Server{
		port:      8080,
		service:   svc,
		validator: validator.New(),
	}
}

func testPopulateMockRepo(t *testing.T, repo *mock_repository.TestRepositoryMock) {
	t.Helper()
	t.Log("Populating Mock Repository")
	validMerchant := model.Merchant{
		Location: model.Location{
			Lat:  6.1753,
			Long: 106.8271,
		},
		Name:             "Store A",
		MerchantCategory: "TODO",
		ImageUrl:         "http://minio",
		ID:               1,
		CreatedAt:        time.Now(),
	}

	tooFarMerchant := model.Merchant{
		Location: model.Location{
			Lat:  35.6764,
			Long: 139.6500,
		},
		Name:             "Store A",
		MerchantCategory: "TODO",
		ImageUrl:         "http://minio",
		ID:               1,
		CreatedAt:        time.Now(),
	}

	validItem := model.Item{
		ID:    1,
		Price: 10,
	}

	// Generate merchants with different locations and IDs
	merchants := []model.Merchant{
		{ID: 1, Name: "Warung Sate Pak Budi", Location: model.Location{Lat: 6.1753, Long: 106.8271}, MerchantCategory: "Street Food", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 2, Name: "Bakso Malang Enak", Location: model.Location{Lat: 6.1760, Long: 106.8280}, MerchantCategory: "Restaurant", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 3, Name: "Nasi Gudeg Bu Sri", Location: model.Location{Lat: 6.1745, Long: 106.8265}, MerchantCategory: "Traditional Food", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 4, Name: "Cafe Kopi Hitam", Location: model.Location{Lat: 6.1770, Long: 106.8290}, MerchantCategory: "Cafe", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 5, Name: "Mie Ayam Pak Tarno", Location: model.Location{Lat: 6.1740, Long: 106.8260}, MerchantCategory: "Street Food", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 6, Name: "Seafood Bu Inem", Location: model.Location{Lat: 6.1780, Long: 106.8300}, MerchantCategory: "Seafood", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 7, Name: "Pizza Corner", Location: model.Location{Lat: 6.1765, Long: 106.8275}, MerchantCategory: "Fast Food", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 8, Name: "Sushi Zen", Location: model.Location{Lat: 6.1755, Long: 106.8285}, MerchantCategory: "Japanese", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 9, Name: "Burger Joint", Location: model.Location{Lat: 6.1750, Long: 106.8270}, MerchantCategory: "Fast Food", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 10, Name: "Taco Fiesta", Location: model.Location{Lat: 6.1775, Long: 106.8295}, MerchantCategory: "Mexican", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 11, Name: "Dim Sum Palace", Location: model.Location{Lat: 6.1748, Long: 106.8268}, MerchantCategory: "Chinese", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 12, Name: "Pasta Italia", Location: model.Location{Lat: 6.1772, Long: 106.8292}, MerchantCategory: "Italian", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 13, Name: "Roti Bakar 88", Location: model.Location{Lat: 6.1742, Long: 106.8262}, MerchantCategory: "Bakery", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 14, Name: "Ayam Geprek Bensu", Location: model.Location{Lat: 6.1778, Long: 106.8298}, MerchantCategory: "Indonesian", ImageUrl: "http://minio", CreatedAt: time.Now()},
		{ID: 15, Name: "Smoothie Bar", Location: model.Location{Lat: 6.1758, Long: 106.8278}, MerchantCategory: "Beverages", ImageUrl: "http://minio", CreatedAt: time.Now()},
	}

	// Generate items for each merchant
	itemSets := [][]model.Item{
		{{ID: 1, Price: 12000, Name: "Sate Ayam"}, {ID: 2, Price: 15000, Name: "Sate Kambing"}},
		{{ID: 3, Price: 8000, Name: "Bakso Urat"}, {ID: 4, Price: 10000, Name: "Bakso Jumbo"}},
		{{ID: 5, Price: 18000, Name: "Gudeg Komplit"}, {ID: 6, Price: 12000, Name: "Gudeg Ayam"}},
		{{ID: 7, Price: 25000, Name: "Cappuccino"}, {ID: 8, Price: 30000, Name: "Latte"}},
		{{ID: 9, Price: 9000, Name: "Mie Ayam Biasa"}, {ID: 10, Price: 12000, Name: "Mie Ayam Bakso"}},
		{{ID: 11, Price: 35000, Name: "Udang Bakar"}, {ID: 12, Price: 40000, Name: "Ikan Gurame"}},
		{{ID: 13, Price: 45000, Name: "Pizza Margherita"}, {ID: 14, Price: 55000, Name: "Pizza Pepperoni"}},
		{{ID: 15, Price: 38000, Name: "Salmon Roll"}, {ID: 16, Price: 42000, Name: "Tuna Sashimi"}},
		{{ID: 17, Price: 22000, Name: "Cheese Burger"}, {ID: 18, Price: 25000, Name: "Beef Burger"}},
		{{ID: 19, Price: 28000, Name: "Beef Tacos"}, {ID: 20, Price: 24000, Name: "Chicken Tacos"}},
		{{ID: 21, Price: 32000, Name: "Har Gow"}, {ID: 22, Price: 30000, Name: "Siu Mai"}},
		{{ID: 23, Price: 48000, Name: "Spaghetti Carbonara"}, {ID: 24, Price: 52000, Name: "Fettuccine Alfredo"}},
		{{ID: 25, Price: 15000, Name: "Roti Bakar Coklat"}, {ID: 26, Price: 18000, Name: "Roti Bakar Keju"}},
		{{ID: 27, Price: 16000, Name: "Ayam Geprek Sambal Ijo"}, {ID: 28, Price: 14000, Name: "Ayam Geprek Original"}},
		{{ID: 29, Price: 20000, Name: "Mango Smoothie"}, {ID: 30, Price: 22000, Name: "Berry Smoothie"}},
	}

	validMerchantWithItems := make([]model.MerchantItem, len(merchants))
	for i, merchant := range merchants {
		validMerchantWithItems[i] = model.MerchantItem{
			Merchant: merchant,
			Items:    itemSets[i],
		}
	}

	repo.Mock.On("GetMerchantByID", mock.Anything, int64(1)).Return(validMerchant, nil)
	repo.Mock.On("GetMerchantByID", mock.Anything, int64(2)).Return(validMerchant, nil)
	repo.Mock.On("GetMerchantByID", mock.Anything, int64(99)).Return(tooFarMerchant, nil)
	repo.Mock.On("GetMerchantByID", mock.Anything, int64(100)).Return(model.Merchant{}, constants.ErrMerchantNotFound)

	repo.Mock.On("GetMerchantByCellID", mock.Anything, mock.Anything).Return(validMerchant, nil)
	repo.Mock.On("ListMerchantWithItems", mock.Anything, mock.Anything).Return(validMerchantWithItems, nil)

	repo.Mock.On("GetItemByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
		return id >= 1 && id < 99
	})).Return(validItem, nil)
	repo.Mock.On("GetItemByID", mock.Anything, int64(100)).Return(model.Item{}, constants.ErrItemNotFound)
}

func testPopulateMockLocationService(t *testing.T, svc *mocklocationservice.MockLocationService) {
	t.Helper()

	neigbhors := []model.Cell{
		{CellID: 1, Resolution: 8}, {CellID: 2, Resolution: 8}, {CellID: 3, Resolution: 8},
		{CellID: 4, Resolution: 8}, {CellID: 5, Resolution: 8}, {CellID: 6, Resolution: 8},
	}

	svc.Mock.On("EstimateDeliveryTimeInMinutes", mock.Anything, mock.Anything).Return(int64(25), nil)
	svc.Mock.On("FindNearby", mock.Anything, mock.Anything, mock.Anything).Return(neigbhors, nil)
	svc.Mock.On("FindCellIDByResolution", mock.Anything, mock.Anything, mock.Anything).Return(model.Cell{CellID: 1, Resolution: 8}, nil)
	svc.Mock.On("FindKRingCellIDs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(neigbhors, nil)
}

// utils.CalculateDistance calculates distance between two points using Haversine formula

func validateMerchantsOrderedByDistance(t *testing.T, userLat, userLon float64, merchants []MerchantWithItem) {
	t.Helper()

	if len(merchants) <= 1 {
		return
	}

	var prevDistance float64 = -1
	for i, merchant := range merchants {
		distance := utils.CalculateDistance(userLat, userLon, merchant.Merchant.Location.Lat, merchant.Merchant.Location.Long)

		if prevDistance >= 0 {
			assert.LessOrEqual(t, prevDistance, distance,
				"Merchant %d (%s) at distance %.2fm should not come before merchant %d (%s) at distance %.2fm",
				i-1, merchants[i-1].Merchant.Name, prevDistance,
				i, merchant.Merchant.Name, distance)
		}

		prevDistance = distance
		t.Logf("Merchant %d: %s - Distance: %.2f meters", i, merchant.Merchant.Name, distance)
	}
}

func TestEstimateOrderPrice(t *testing.T) {
	t.Skip()
	s := testPurchaseSetup(t)

	validReq := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "2",
				IsStartingPoint: false,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	invalidStartingPointReq := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "2",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	invaliMerchantTooFarReq := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "99",
				IsStartingPoint: false,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	InvalidMerchantNotFound := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "100",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "1",
				IsStartingPoint: false,
				Items: []OrderItemRequest{
					{ItemID: "3", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	InvalidItemNotFound := OrderEstimationRequest{
		UserLocation: LocationRequest{
			Lat:  6.1674,
			Long: 106.8209,
		},

		Orders: []OrderRequest{
			{
				MerchantID:      "1",
				IsStartingPoint: true,
				Items: []OrderItemRequest{
					{ItemID: "1", Quantity: 1},
					{ItemID: "2", Quantity: 2},
				},
			},
			{
				MerchantID:      "2",
				IsStartingPoint: false,
				Items: []OrderItemRequest{
					{ItemID: "100", Quantity: 1},
					{ItemID: "4", Quantity: 2},
				},
			},
		},
	}

	t.Run("Valid", func(t *testing.T) {
		reqBody, err := json.Marshal(validReq)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.NotEqual(t, 0, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})

	t.Run("Invalid_StartingPoint", func(t *testing.T) {
		reqBody, err := json.Marshal(invalidStartingPointReq)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})

	t.Run("Invalid_MerchantTooFar", func(t *testing.T) {
		reqBody, err := json.Marshal(invaliMerchantTooFarReq)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})

	t.Run("Invalid_MerchantNotFound", func(t *testing.T) {
		reqBody, err := json.Marshal(InvalidMerchantNotFound)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.Equal(t, 404, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})

	t.Run("Invalid_ItemNotFound", func(t *testing.T) {
		reqBody, err := json.Marshal(InvalidItemNotFound)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/users/estimate", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		s.EstimateOrderPrice(w, req)

		resp := w.Result()
		assert.Equal(t, 404, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})
}

func TestGetNearbyMerchants(t *testing.T) {
	s := testPurchaseSetup(t)

	userLocation := LocationRequest{
		Lat:  6.1674,
		Long: 106.8209,
	}

	t.Run("Valid_DistanceOrdering", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		// Parse response and validate distance ordering
		var response FindNearbyMerchantsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		for _, data := range response.Data {
			t.Log(data)
		}

		t.Logf("%+v", response.Meta)

		validateMerchantsOrderedByDistance(t, userLocation.Lat, userLocation.Long, response.Data)
		assert.NotEmpty(t, response.Data, "Should return merchants")
	})

	t.Run("ValidWithFilter_MerchantID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby?merchantId=2", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		// Parse response and validate distance ordering
		var response FindNearbyMerchantsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		for _, data := range response.Data {
			t.Log(data)
		}

		validateMerchantsOrderedByDistance(t, userLocation.Lat, userLocation.Long, response.Data)
		assert.NotEmpty(t, response.Data, "Should return merchants")
	})

	t.Run("ValidWithFilter_Name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby?name=bat", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		// Parse response and validate distance ordering
		var response FindNearbyMerchantsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		for _, data := range response.Data {
			t.Log(data)
		}

		validateMerchantsOrderedByDistance(t, userLocation.Lat, userLocation.Long, response.Data)
		assert.NotEmpty(t, response.Data, "Should return merchants")
	})

	t.Run("ValidWithPagination_DistanceOrdering", func(t *testing.T) {
		// t.Skip()
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby?limit=3&offset=1", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		// Parse response and validate distance ordering
		var response FindNearbyMerchantsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		validateMerchantsOrderedByDistance(t, userLocation.Lat, userLocation.Long, response.Data)
		assert.LessOrEqual(t, len(response.Data), 3, "Should respect limit parameter")
	})

	t.Run("EdgeCase_OffsetGreaterThanTotal", func(t *testing.T) {
		// t.Skip()
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby?limit=10&offset=100", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		// Parse response - should return all merchants when offset exceeds total
		var response FindNearbyMerchantsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		validateMerchantsOrderedByDistance(t, userLocation.Lat, userLocation.Long, response.Data)
		assert.NotEmpty(t, response.Data, "Should return all merchants when offset > total")
	})

	t.Run("EdgeCase_LimitExceedsRemaining", func(t *testing.T) {
		// t.Skip()
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby?limit=100&offset=3", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		// Parse response and validate distance ordering
		var response FindNearbyMerchantsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		validateMerchantsOrderedByDistance(t, userLocation.Lat, userLocation.Long, response.Data)
	})

	t.Run("InvalidCoordinate", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby", nil)
		req.SetPathValue("coordinate", "invalid")
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/merchants/nearby", nil)
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 405, resp.StatusCode)
	})
}

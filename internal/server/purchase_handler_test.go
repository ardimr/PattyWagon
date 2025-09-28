package server

import (
	"PattyWagon/internal/constants"
	imagecompressor "PattyWagon/internal/image_compressor"
	mocklocationservice "PattyWagon/internal/mock_location_service"
	"PattyWagon/internal/mock_repository"
	"PattyWagon/internal/model"
	"PattyWagon/internal/service"
	"PattyWagon/internal/storage"
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
	repo := &mock_repository.TestRepositoryMock{}
	storage := storage.New("localhost:9000", "team-solid", "@team-solid", storage.Option{MaxConcurrent: 5})
	imageCompressor := imagecompressor.New(5, 50)
	locationSvc := &mocklocationservice.MockLocationService{}
	svc := service.New(repo, storage, imageCompressor, locationSvc)

	testPopulateMockRepo(t, repo)
	testPopulateMockLocationService(t, locationSvc)
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

	validItems := []model.Item{
		{ID: 1, Price: 10, Name: "Bakso"},
		{ID: 2, Price: 15, Name: "Mie Ayam"},
	}
	validMerchantWithItems := []model.MerchantItem{
		{Merchant: validMerchant, Items: validItems},
		{Merchant: validMerchant, Items: validItems},
		{Merchant: validMerchant, Items: validItems},
		{Merchant: validMerchant, Items: validItems},
		{Merchant: validMerchant, Items: validItems},
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
		{ID: 1, CellID: 1, MerchantID: 1}, {ID: 2, CellID: 2, MerchantID: 2}, {ID: 3, CellID: 3, MerchantID: 3},
		{ID: 4, CellID: 4, MerchantID: 1}, {ID: 5, CellID: 5, MerchantID: 2}, {ID: 6, CellID: 6, MerchantID: 3},
	}

	svc.Mock.On("EstimateDeliveryTimeInMinutes", mock.Anything, mock.Anything).Return(int64(25), nil)
	svc.Mock.On("FindNearby", mock.Anything, mock.Anything, mock.Anything).Return(neigbhors, nil)
}

func TestEstimateOrderPrice(t *testing.T) {
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

	filter := FindNearbyMerchantRequest{}

	userLocation := LocationRequest{
		Lat:  6.1674,
		Long: 106.8209,
	}

	t.Run("Valid", func(t *testing.T) {
		reqBody, err := json.Marshal(filter)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodGet, "/v1/merchants/nearby", bytes.NewBuffer(reqBody))
		req.SetPathValue("coordinate", fmt.Sprintf("%f,%f", userLocation.Lat, userLocation.Long))
		w := httptest.NewRecorder()

		s.FindNearbyMerchants(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)
		fmt.Printf("Response Body: %v\n", w.Body.String())
	})
}

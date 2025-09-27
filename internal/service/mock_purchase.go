package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"io"

	"github.com/kwahome/go-haversine/pkg/haversine"
)

type MockPurchaseService struct {
	repository      Repository
	storage         Storage
	imageCompressor ImageCompressor
}

func NewMockPurchaseService(repository Repository, storage Storage, imageCompressor ImageCompressor) *MockPurchaseService {
	return &MockPurchaseService{
		repository:      repository,
		storage:         storage,
		imageCompressor: imageCompressor,
	}
}

func (mock *MockPurchaseService) EmailLogin(ctx context.Context, user string, password string) (string, string, error) {
	return "", "", nil
}

func (mock *MockPurchaseService) Register(ctx context.Context, user model.User, password string) (string, error) {
	return "", nil
}

func (mock *MockPurchaseService) IsUserExist(ctx context.Context, userID int64) (bool, error) {
	return true, nil
}

func (mock *MockPurchaseService) UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (model.File, error) {
	return model.File{}, nil
}

// Purchase
func (mock *MockPurchaseService) EstimateOrderPrice(ctx context.Context, req model.OrderEstimation) (model.EstimationPrice, error) {
	isStartingPointValid := false

	for _, order := range req.Orders {
		isStartingPointValid = isStartingPointValid != order.IsStartingPoint

		var merchantLoc model.Location
		if order.MerchantID == 99 {
			merchantLoc = model.Location{
				Lat:  35.6764,
				Long: 139.6500,
			}
		} else {
			merchantLoc = model.Location{
				Lat:  6.1753,
				Long: 106.8271,
			}
		}
		if !mock.validateDistance(ctx, req.UserLocation, merchantLoc) {
			return model.EstimationPrice{}, constants.ErrMerchantTooFar
		}
	}

	if !isStartingPointValid {
		return model.EstimationPrice{}, constants.ErrInvalidStartingPoint
	}

	return model.EstimationPrice{
		ID:                         1,
		EstimatedDeliveryInMinutes: 5,
		TotalPrice:                 10,
	}, nil
}

func (mock *MockPurchaseService) validateDistance(ctx context.Context, user, merchant model.Location) bool {
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

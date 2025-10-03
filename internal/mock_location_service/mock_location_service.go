package mocklocationservice

import (
	"PattyWagon/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationService) EstimateDeliveryTimeInMinutes(ctx context.Context, locations []model.Location) (int64, error) {
	args := m.Called(ctx, locations)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLocationService) FindNearby(ctx context.Context, location model.Location, searchingLevel int) ([]model.Cell, error) {
	args := m.Called(ctx, location, searchingLevel)
	return args.Get(0).([]model.Cell), args.Error(1)
}

func (m *MockLocationService) GetAllCellIDs(ctx context.Context, location model.Location) ([]model.Cell, error) {
	args := m.Called(ctx, location, location)
	return args.Get(0).([]model.Cell), args.Error(1)
}

func (m *MockLocationService) FindCellIDByResolution(ctx context.Context, location model.Location, resolution int) (model.Cell, error) {
	args := m.Called(ctx, location, resolution)
	return args.Get(0).(model.Cell), args.Error(1)
}

func (m *MockLocationService) FindKRingCellIDs(ctx context.Context, location model.Location, resolution, k int) ([]model.Cell, error) {
	args := m.Called(ctx, location, resolution, k)
	return args.Get(0).([]model.Cell), args.Error(1)
}

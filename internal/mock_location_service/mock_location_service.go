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

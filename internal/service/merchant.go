package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
)

func (s *Service) GetMerchant(ctx context.Context, merchantID int64) (model.Merchant, error) {
	merchant, err := s.repository.GetMerchantByID(ctx, merchantID)
	if err != nil {
		return model.Merchant{}, constants.WrapError(constants.ErrFailedToGetMerchant, err)
	}
	return merchant, nil
}
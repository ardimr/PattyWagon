package service

import (
	"PattyWagon/internal/model"
	"context"
)

func (s *Service) CreateMerchant(ctx context.Context, req model.Merchant) (res int64, err error) {
	//
	// Insert New Merchant
	//
	newMerchant := model.Merchant{
		UserID:    req.UserID,
		Name:      req.Name,
		Category:  req.Category,
		ImageURL:  req.ImageURL,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	res, err = s.repository.InsertMerchant(ctx, newMerchant)
	if err != nil {
		return 0, err
	}

	return res, nil
}

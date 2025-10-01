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

func (s *Service) GetMerchants(ctx context.Context, req model.FilterMerchant) (res []model.Merchant, err error) {
	//
	// Get Merchants
	//
	paramsFetchMerchant := model.FilterMerchant{
		MerchantID:       req.MerchantID,
		Limit:            req.Limit,
		Offset:           req.Offset,
		Name:             req.Name,
		MerchantCategory: req.MerchantCategory,
		CreatedAt:        req.CreatedAt,
	}
	res, err = s.repository.GetMerchants(ctx, paramsFetchMerchant)
	if err != nil {
		return []model.Merchant{}, err
	}

	return res, nil
}

package service

import (
	"PattyWagon/internal/model"
	"context"
)

func (s *Service) CreateItems(ctx context.Context, req model.Item) (res int64, err error) {
	//
	// Check Existing Merchant
	//
	merchantID := req.MerchantID
	_, err = s.repository.MerchantExists(ctx, merchantID)
	if err != nil {
		return 0, err
	}
	//
	// Insert New Items
	//
	newItem := model.Item{
		MerchantID: req.MerchantID,
		Name:       req.Name,
		Category:   req.Category,
		Price:      req.Price,
		ImageURL:   req.ImageURL,
	}
	res, err = s.repository.CreateItems(ctx, newItem)
	if err != nil {
		return 0, err
	}

	return res, nil
}

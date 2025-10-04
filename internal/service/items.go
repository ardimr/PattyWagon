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

func (s *Service) GetItems(ctx context.Context, req model.FilterItem) (res []model.Item, err error) {
	//
	// Get Items
	//
	paramsFetchItem := model.FilterItem{
		ItemID:          req.ItemID,
		Name:            req.Name,
		ProductCategory: req.ProductCategory,
		Limit:           req.Limit,
		Offset:          req.Offset,
		CreatedAt:       req.CreatedAt,
	}
	res, err = s.repository.GetItems(ctx, paramsFetchItem)
	if err != nil {
		return []model.Item{}, err
	}

	return
}

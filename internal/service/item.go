package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
)

func (s *Service) GetItem(ctx context.Context, itemID int64) (model.Item, error) {
	item, err := s.repository.GetItemByID(ctx, itemID)
	if err != nil {
		return model.Item{}, constants.WrapError(constants.ErrFailedToGetItem, err)
	}
	return item, nil
}
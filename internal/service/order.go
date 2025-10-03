package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"errors"
)

func (s *Service) GetOrder(ctx context.Context, orderID int64) (model.Order, error) {
	order, err := s.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (s *Service) GetOrderDetail(ctx context.Context, orderDetailID int64) (model.OrderDetail, error) {
	orderDetail, err := s.repository.GetOrderDetailByID(ctx, orderDetailID)
	if err != nil {
		return model.OrderDetail{}, err
	}
	return orderDetail, nil
}

func (s *Service) AddItemToCart(ctx context.Context, userID, merchantID, itemID int64, quantity int32) (model.OrderDetail, error) {
	// Start a database transaction to ensure atomicity
	tx, err := s.repository.BeginTx(ctx, nil)
	if err != nil {
		return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToBeginTransaction, err)
	}

	txRepo := s.repository.WithTx(tx).(TxRepository)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	merchant, err := s.GetMerchant(ctx, merchantID)
	if err != nil {
		return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetMerchant, err)
	}

	item, err := s.GetItem(ctx, itemID)
	if err != nil {
		return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetItem, err)
	}

	order, err := txRepo.GetUnpurchasedOrderByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, constants.ErrNoUnpurchasedOrder) {
			order, err = txRepo.CreateOrder(ctx, userID, 0, false)
			if err != nil {
				return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCreateUnpurchasedOrder, err)
			}
		} else {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetUnpurchasedOrder, err)
		}
	}

	orderDetail, err := txRepo.GetOrderDetailByOrderIDAndMerchantID(ctx, order.ID, merchantID)
	if err != nil {
		if errors.Is(err, constants.ErrOrderDetailNotFound) {
			orderDetail, err = txRepo.CreateOrderDetail(ctx, order.ID, merchantID, merchant.Name, merchant.Category.String, merchant.ImageURL, merchant.Latitude, merchant.Longitude)
			if err != nil {
				return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCreateOrderDetail, err)
			}
		} else {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetOrderDetail, err)
		}
	}

	orderItem, err := txRepo.GetOrderItemByOrderDetailIDAndItemID(ctx, orderDetail.ID, itemID)
	if err != nil {
		if errors.Is(err, constants.ErrOrderItemNotFound) {
			pricePerItem := int64(0)
			if item.Price.Valid {
				pricePerItem = item.Price.Int64
			}

			_, err = txRepo.CreateOrderItem(ctx, orderDetail.ID, itemID, item.Name, item.Category.String, item.FileURI.String, pricePerItem, quantity, pricePerItem*int64(quantity))
			if err != nil {
				return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCreateOrderItem, err)
			}
		} else {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetOrderItem, err)
		}
	} else {
		newQuantity := orderItem.Quantity + quantity
		newTotalPrice := orderItem.PricePerItem * int64(newQuantity)
		_, err = txRepo.UpdateOrderItem(ctx, orderItem.ID, orderItem.OrderDetailID, orderItem.ItemName, orderItem.ProductCategory, orderItem.ItemImageURL, orderItem.PricePerItem, newQuantity, newTotalPrice)
		if err != nil {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToUpdateOrderItem, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCommitTransaction, err)
	}

	return orderDetail, nil
}

func (s *Service) CreateUnpurchasedOrder(ctx context.Context, userID int64) (model.Order, error) {
	order, err := s.repository.CreateOrder(ctx, userID, 0, false)
	if err != nil {
		return model.Order{}, constants.WrapError(constants.ErrFailedToCreateOrder, err)
	}

	return order, nil
}

func (s *Service) GetPurchasedOrdersPagination(ctx context.Context, userID int64, limit, offset int) ([]model.Order, error) {
	orders, err := s.repository.GetOrdersByUserIDAndPurchased(ctx, userID, true, limit, offset)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) UpdateOrderToPurchased(ctx context.Context, orderID int64) error {
	order, err := s.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		return constants.WrapError(constants.ErrFailedToGetOrder, err)
	}

	_, err = s.repository.UpdateOrder(ctx, orderID, order.OrderEstimationID, true)
	if err != nil {
		return constants.WrapError(constants.ErrFailedToUpdateOrder, err)
	}

	return nil
}

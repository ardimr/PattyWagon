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

	order, err := s.repository.GetUnpurchasedOrderByUserIDWithTx(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, constants.ErrNoUnpurchasedOrder) {
			order, err = s.repository.CreateOrderWithTx(ctx, tx, userID, 0, false)
			if err != nil {
				return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCreateUnpurchasedOrder, err)
			}
		} else {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetUnpurchasedOrder, err)
		}
	}

	orderDetail, err := s.repository.GetOrderDetailByOrderIDAndMerchantIDWithTx(ctx, tx, order.ID, merchantID)
	if err != nil {
		if errors.Is(err, constants.ErrOrderDetailNotFound) {
			orderDetail, err = s.repository.CreateOrderDetailWithTx(ctx, tx, order.ID, merchantID, merchant.Name, merchant.Category, merchant.ImageURL, merchant.Latitude, merchant.Longitude)
			if err != nil {
				return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCreateOrderDetail, err)
			}
		} else {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetOrderDetail, err)
		}
	}

	orderItem, err := s.repository.GetOrderItemByOrderDetailIDAndItemIDWithTx(ctx, tx, orderDetail.ID, itemID)
	if err != nil {
		if errors.Is(err, constants.ErrOrderItemNotFound) {
			pricePerItem := item.Price
			// if item.Price.Valid {
			// 	pricePerItem = item.Price.Int64
			// }

			_, err = s.repository.CreateOrderItemWithTx(ctx, tx, orderDetail.ID, itemID, item.Name, item.Category, item.ImageURL, pricePerItem, quantity, pricePerItem*float64(quantity))
			if err != nil {
				return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToCreateOrderItem, err)
			}
		} else {
			return model.OrderDetail{}, constants.WrapError(constants.ErrFailedToGetOrderItem, err)
		}
	} else {
		newQuantity := orderItem.Quantity + quantity
		newTotalPrice := orderItem.PricePerItem * float64(newQuantity)
		_, err = s.repository.UpdateOrderItemWithTx(ctx, tx, orderItem.ID, orderItem.OrderDetailID, orderItem.ItemName, orderItem.ProductCategory, orderItem.ItemImageURL, orderItem.PricePerItem, newQuantity, newTotalPrice)
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

package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"context"
	"sync"
)

func (s *Service) CreateOrderEstimation(ctx context.Context, req model.EstimationRequest) (model.EstimationResponse, error) {
	// TODO: Implement create order estimation
	return model.EstimationResponse{}, nil
}

func (s *Service) GetMerchantItemDataConcurrently(ctx context.Context, items []model.EstimationRequestItem) ([]model.MerchantItemData, error) {
	if len(items) == 0 {
		return []model.MerchantItemData{}, nil
	}

	type result struct {
		data  model.MerchantItemData
		err   error
		index int
	}

	resultChan := make(chan result, len(items))

	var wg sync.WaitGroup
	for i, item := range items {
		wg.Add(1)
		go func(index int, reqItem model.EstimationRequestItem) {
			defer wg.Done()

			// Get merchant and item concurrently
			merchantChan := make(chan struct {
				merchant model.Merchant
				err      error
			}, 1)
			itemChan := make(chan struct {
				item model.Item
				err  error
			}, 1)

			go func() {
				merchant, err := s.GetMerchant(ctx, reqItem.MerchantID)
				merchantChan <- struct {
					merchant model.Merchant
					err      error
				}{merchant, err}
			}()

			go func() {
				item, err := s.GetItem(ctx, reqItem.ItemID)
				itemChan <- struct {
					item model.Item
					err  error
				}{item, err}
			}()

			merchantResult := <-merchantChan
			itemResult := <-itemChan

			if merchantResult.err != nil {
				resultChan <- result{
					err:   constants.WrapError(constants.ErrFailedToGetMerchant, merchantResult.err),
					index: index,
				}
				return
			}

			if itemResult.err != nil {
				resultChan <- result{
					err:   constants.WrapError(constants.ErrFailedToGetItem, itemResult.err),
					index: index,
				}
				return
			}

			resultChan <- result{
				data: model.MerchantItemData{
					Merchant: merchantResult.merchant,
					Item:     itemResult.item,
					Quantity: reqItem.Quantity,
				},
				index: index,
			}
		}(i, item)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	merchantItemData := make([]model.MerchantItemData, len(items))
	for result := range resultChan {
		if result.err != nil {
			return nil, result.err
		}
		merchantItemData[result.index] = result.data
	}

	return merchantItemData, nil
}

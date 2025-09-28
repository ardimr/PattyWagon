package repository

import (
	"PattyWagon/internal/model"
	"context"
)

func (q *Queries) GetMerchantByID(ctx context.Context, id int64) (model.Merchant, error) {
	panic("unimplemented")
}

func (q *Queries) GetMerchantByCellID(ctx context.Context, cellID int64) (model.Merchant, error) {
	panic("unimplemented")
}

package model

type ListMerchantWithItemParams struct {
	Cell *Cell
	MerchantParams
}

type MerchantParams struct {
	MerchantID       *int64
	Limit            int
	Offset           int
	Name             *string
	MerchantCategory *string
	SortingOrder     *string
}

type MerchantCategory string

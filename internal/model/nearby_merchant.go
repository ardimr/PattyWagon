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

type MerchantItem struct {
	Merchant Merchant
	Items    []Item
}

type MerchantCategory string

type FindNerbyMerchantParams struct {
	UserLocation Location
	MerchantParams
}

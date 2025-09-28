package model

type Merchant struct {
	Location Location
}

type Cell struct {
	ID         int64
	CellID     int64
	MerchantID int64
}

type MerchantCategory string

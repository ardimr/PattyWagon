package service

import (
	"PattyWagon/internal/model"
	"context"
	"database/sql"
)

type Service struct {
	repository      Repository
	storage         Storage
	imageCompressor ImageCompressor
	locationService LocationService
	merchantCounter MerchantCounter
}

type Repository interface {
	// Transaction methods
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// User Repository
	InsertUser(ctx context.Context, user model.User, passwordHash string) (model.User, error)
	SelectUserCredentialsByUsernameAndRole(ctx context.Context, username string, role int16) (res model.User, err error)

	// File Repository
	GetFileUpload(ctx context.Context, id int64) (model.File, error)
	InsertFile(ctx context.Context, file model.File) (model.File, error)
	GetFileByFileID(ctx context.Context, fileID string) (res model.File, err error)
	FileExists(ctx context.Context, fileID string) (bool, error)

	// Merchant Repository
	InsertMerchant(ctx context.Context, data model.Merchant) (res int64, err error)
	GetMerchants(ctx context.Context, filter model.FilterMerchant) (res []model.Merchant, err error)
	MerchantExists(ctx context.Context, merchantID int64) (res bool, err error)
	BulkInsertMerchantLocations(ctx context.Context, locations []model.MerchantLocation) error

	CreateItems(ctx context.Context, item model.Item) (int64, error)
	GetItems(ctx context.Context, filter model.FilterItem) (res []model.Item, err error)
	// Merchant Repository
	GetMerchantByCellID(ctx context.Context, cellID int64) (model.Merchant, error)
	ListMerchantWithItems(ctx context.Context, params model.ListMerchantWithItemParams) ([]model.MerchantItem, error)
	// Item

	GetMerchantWithItems(ctx context.Context, merchantID int64) (model.MerchantItem, error)
	// Merchant Repository
	GetMerchantByID(ctx context.Context, id int64) (model.Merchant, error)

	// Item Repository
	GetItemByID(ctx context.Context, id int64) (model.Item, error)

	// Order Repository
	CreateOrder(ctx context.Context, userID int64, orderEstimationID int64, isPurchased bool) (model.Order, error)
	CreateOrderWithTx(ctx context.Context, tx *sql.Tx, userID int64, orderEstimationID int64, isPurchased bool) (model.Order, error)
	GetOrderByID(ctx context.Context, id int64) (model.Order, error)
	GetOrderByEstimationID(ctx context.Context, estimationId int64) (model.Order, error)
	GetOrdersByUserIDAndPurchased(ctx context.Context, userID int64, isPurchased bool, limit, offset int) ([]model.Order, error)
	GetUnpurchasedOrderByUserIDWithTx(ctx context.Context, tx *sql.Tx, userID int64) (model.Order, error)
	UpdateOrder(ctx context.Context, id int64, orderEstimationID int64, isPurchased bool) (model.Order, error)

	// Order Detail Repository
	CreateOrderDetailWithTx(ctx context.Context, tx *sql.Tx, orderID, merchantID int64, merchantName, merchantCategory, merchantImageURL string, merchantLatitude, merchantLongitude float64) (model.OrderDetail, error)
	GetOrderDetailByID(ctx context.Context, id int64) (model.OrderDetail, error)
	GetOrderDetailByOrderIDAndMerchantIDWithTx(ctx context.Context, tx *sql.Tx, orderID, merchantID int64) (model.OrderDetail, error)

	// Order Item Repository
	CreateOrderItemWithTx(ctx context.Context, tx *sql.Tx, orderDetailID, itemID int64, itemName, productCategory, itemImageURL string, pricePerItem float64, quantity int32, totalPrice float64) (model.OrderItem, error)
	GetOrderItemByOrderDetailIDAndItemIDWithTx(ctx context.Context, tx *sql.Tx, orderDetailID, itemID int64) (model.OrderItem, error)
	UpdateOrderItemWithTx(ctx context.Context, tx *sql.Tx, id, orderDetailID int64, itemName, productCategory, itemImageURL string, pricePerItem float64, quantity int32, totalPrice float64) (model.OrderItem, error)
}

type Storage interface {
	UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error)
}

type ImageCompressor interface {
	Compress(ctx context.Context, src string) (string, error)
}

type LocationService interface {
	GetAllCellIDs(ctx context.Context, location model.Location) ([]model.Cell, error)
	FindCellIDByResolution(ctx context.Context, location model.Location, resolution int) (model.Cell, error)
	FindKRingCellIDs(ctx context.Context, location model.Location, resolution, k int) ([]model.Cell, error)
}

type MerchantCounter interface {
	Increment()
	Get() int64
}

func New(repository Repository, storage Storage, imageCompressor ImageCompressor, locationService LocationService, merchantCounter MerchantCounter) *Service {
	return &Service{
		repository:      repository,
		storage:         storage,
		imageCompressor: imageCompressor,
		locationService: locationService,
		merchantCounter: merchantCounter,
	}
}

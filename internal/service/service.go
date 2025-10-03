package service

import (
	"PattyWagon/internal/model"
	"context"
	"database/sql"
)

type Service struct {
	repository      Repository
	TxRepository    TxRepository
	storage         Storage
	imageCompressor ImageCompressor
}

type TxRepository interface {
	Repository
}

type Repository interface {
	// Transaction methods
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	WithTx(tx *sql.Tx) any

	// User
	// User Repository
	InsertUser(ctx context.Context, user model.User, passwordHash string) (model.User, error)
	SelectUserCredentialsByEmail(ctx context.Context, phone string) (model.User, error)

	IsUserExist(ctx context.Context, userID int64) (bool, error)

	// File
	GetFileUpload(ctx context.Context, id int64) (model.File, error)

	// File Repository
	InsertFile(ctx context.Context, file model.File) (model.File, error)
	GetFileByFileID(ctx context.Context, fileID string) (res model.File, err error)
	FileExists(ctx context.Context, fileID string) (bool, error)

	// Merchant Repository
	GetMerchantByID(ctx context.Context, id int64) (model.Merchant, error)

	// Item Repository
	GetItemByID(ctx context.Context, id int64) (model.Item, error)

	// Order Repository
	CreateOrder(ctx context.Context, userID int64, orderEstimationID int64, isPurchased bool) (model.Order, error)
	GetOrderByID(ctx context.Context, id int64) (model.Order, error)
	GetOrderByEstimationID(ctx context.Context, estimationId int64) (model.Order, error)
	GetOrdersByUserIDAndPurchased(ctx context.Context, userID int64, isPurchased bool, limit, offset int) ([]model.Order, error)
	GetUnpurchasedOrderByUserID(ctx context.Context, userID int64) (model.Order, error)
	UpdateOrder(ctx context.Context, id int64, orderEstimationID int64, isPurchased bool) (model.Order, error)

	// Order Detail Repository
	CreateOrderDetail(ctx context.Context, orderID, merchantID int64, merchantName, merchantCategory, merchantImageURL string, merchantLatitude, merchantLongitude float64) (model.OrderDetail, error)
	GetOrderDetailByID(ctx context.Context, id int64) (model.OrderDetail, error)
	GetOrderDetailByOrderIDAndMerchantID(ctx context.Context, orderID, merchantID int64) (model.OrderDetail, error)

	// Order Item Repository
	CreateOrderItem(ctx context.Context, orderDetailID, itemID int64, itemName, productCategory, itemImageURL string, pricePerItem int64, quantity int32, totalPrice int64) (model.OrderItem, error)
	GetOrderItemByOrderDetailIDAndItemID(ctx context.Context, orderDetailID, itemID int64) (model.OrderItem, error)
	UpdateOrderItem(ctx context.Context, id, orderDetailID int64, itemName, productCategory, itemImageURL string, pricePerItem int64, quantity int32, totalPrice int64) (model.OrderItem, error)
}

type Storage interface {
	UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error)
}

type ImageCompressor interface {
	Compress(ctx context.Context, src string) (string, error)
}

func New(repository Repository, txRepo TxRepository, storage Storage, imageCompressor ImageCompressor) *Service {
	return &Service{
		repository:      repository,
		TxRepository:    txRepo,
		storage:         storage,
		imageCompressor: imageCompressor,
	}
}

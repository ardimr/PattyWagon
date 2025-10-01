package service

import (
	"PattyWagon/internal/model"
	"context"
)

type Service struct {
	repository      Repository
	storage         Storage
	imageCompressor ImageCompressor
	locationService LocationService
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
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
	GetMerchantByCellID(ctx context.Context, cellID int64) (model.Merchant, error)
	ListMerchantWithItems(ctx context.Context, params model.ListMerchantWithItemParams) ([]model.MerchantItem, error)
	// Item
	GetItemByID(ctx context.Context, id int64) (model.Item, error)
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
	FindKRingCellIDs(ctx context.Context, location model.Location, k int) ([]model.Cell, error)
}

func New(repository Repository, storage Storage, imageCompressor ImageCompressor, locationService LocationService) *Service {
	return &Service{
		repository:      repository,
		storage:         storage,
		imageCompressor: imageCompressor,
		locationService: locationService,
	}
}

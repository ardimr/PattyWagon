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
	SelectUserCredentialsByUsernameAndRole(ctx context.Context, username string, role int16) (res model.User, err error)

	// File
	GetFileUpload(ctx context.Context, id int64) (model.File, error)

	// File Repository
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

func New(repository Repository, storage Storage, imageCompressor ImageCompressor, LocationService LocationService) *Service {
	return &Service{
		repository:      repository,
		storage:         storage,
		imageCompressor: imageCompressor,
		locationService: LocationService,
	}
}

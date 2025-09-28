package mock_repository

import (
	"PattyWagon/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type TestRepositoryMock struct {
	mock.Mock
}

func (r *TestRepositoryMock) InsertUser(ctx context.Context, user model.User, passwordHash string) (model.User, error) {
	args := r.Called(ctx, user, passwordHash)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *TestRepositoryMock) SelectUserCredentialsByEmail(ctx context.Context, phone string) (model.User, error) {
	args := r.Called(ctx, phone)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *TestRepositoryMock) IsUserExist(ctx context.Context, userID int64) (bool, error) {
	args := r.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (r *TestRepositoryMock) GetFileUpload(ctx context.Context, id int64) (model.File, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(model.File), args.Error(1)
}

func (r *TestRepositoryMock) InsertFile(ctx context.Context, file model.File) (model.File, error) {
	args := r.Called(ctx, file)
	return args.Get(0).(model.File), args.Error(1)
}

func (r *TestRepositoryMock) GetFileByFileID(ctx context.Context, fileID string) (model.File, error) {
	args := r.Called(ctx, fileID)
	return args.Get(0).(model.File), args.Error(1)
}

func (r *TestRepositoryMock) FileExists(ctx context.Context, fileID string) (bool, error) {
	args := r.Called(ctx, fileID)
	return args.Bool(0), args.Error(1)
}

func (r *TestRepositoryMock) GetMerchantByID(ctx context.Context, id int64) (model.Merchant, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(model.Merchant), args.Error(1)
}

func (r *TestRepositoryMock) GetItemByID(ctx context.Context, id int64) (model.Item, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(model.Item), args.Error(1)
}

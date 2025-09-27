package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestMinioStorage(t *testing.T) *MinioStorage {
	t.Helper()
	s3Endpoint := "localhost:9000"
	s3AccessKeyID := "team-solid"
	s3SecretAccessKey := "@team-solid"
	return New(s3Endpoint, s3AccessKeyID, s3SecretAccessKey, Option{
		MaxConcurrent: 5,
	})
}

func TestStorage(t *testing.T) {
	minioStorage := setupTestMinioStorage(t)

	t.Run("Upload_FromFile", func(t *testing.T) {
		bucket := "images"
		localPath := "testdata/sample.jpg"
		remotePath := "sample.jpg"

		url, err := minioStorage.UploadFile(context.TODO(), bucket, localPath, remotePath)
		assert.Nil(t, err)
		assert.NotEmpty(t, url)
	})
}

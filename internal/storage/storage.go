package storage

import (
	"PattyWagon/observability"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	S3AccessKeyID            = os.Getenv("S3_ACCESS_KEY_ID")
	S3SecretAccessKey        = os.Getenv("S3_SECRET_ACCESS_KEY")
	S3Endpoint               = os.Getenv("S3_ENDPOINT")
	S3Bucket                 = os.Getenv("S3_BUCKET")
	S3MaxConcurrentUpload, _ = strconv.Atoi(os.Getenv("S3_MAX_CONCURRENT_UPLOAD"))
)

type MinioStorage struct {
	semaphore         chan struct{}
	client            *minio.Client
	s3Endpoint        string
	s3AccessKeyID     string
	s3SecretAccessKey string
	secure            bool
}

type Option struct {
	MaxConcurrent int64
}

func New(
	s3Endpoint, s3AccessKeyID, s3SecretAccessKey string, option Option,
) *MinioStorage {

	mc, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyID, s3SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &MinioStorage{
		semaphore:         make(chan struct{}, option.MaxConcurrent),
		client:            mc,
		s3Endpoint:        s3Endpoint,
		s3AccessKeyID:     s3AccessKeyID,
		s3SecretAccessKey: s3SecretAccessKey,
		secure:            false,
	}
}

func (s *MinioStorage) UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error) {
	ctx, span := observability.Tracer.Start(ctx, "storage.s3_upload")
	defer span.End()

	select {
	case s.semaphore <- struct{}{}:
		defer func() {
			<-s.semaphore
		}()
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("upload queue timeout")
	}

	_, err := s.client.FPutObject(ctx, bucket, remotePath, localPath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("http://%s/%s/%s", s.s3Endpoint, bucket, remotePath)
	return uri, nil
}

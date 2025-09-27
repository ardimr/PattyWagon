package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/storage"
	"WeenieHut/internal/utils"
	"WeenieHut/observability"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (s *Service) UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (model.File, error) {
	ctx, span := observability.Tracer.Start(ctx, "service.file_upload")
	defer span.End()

	var result model.File

	if sizeInBytes > constants.MaxUploadSizeInBytes {
		return result, constants.ErrMaximumFileSize
	}
	if err := utils.ValidateFileExtensions(filename, constants.AllowedExtensions); err != nil {
		return result, constants.ErrInvalidFileType
	}

	bucket := storage.S3Bucket
	identifier := uuid.NewString()

	tempFilepath := filepath.Join("/tmp", fmt.Sprintf("%s_%s", identifier, filename))
	tempFile, err := os.OpenFile(tempFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return result, err
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	n, err := io.Copy(tempFile, file)
	if err != nil {
		return result, err
	}

	log.Printf("written size: %d filename: %s", n, tempFile.Name())

	// Compress
	thumbnailPath, err := s.imageCompressor.Compress(ctx, tempFile.Name())
	if err != nil {
		return result, err
	}
	defer os.Remove(thumbnailPath)

	// Upload to object storage
	remotePath := fmt.Sprintf("%s_%s", identifier, filename)
	uri, err := s.storage.UploadFile(ctx, bucket, tempFile.Name(), remotePath)
	if err != nil {
		return result, fmt.Errorf("error uploading original file: %w", err)
	}

	thumbnailName := filepath.Base(thumbnailPath)
	thumbnailUri, err := s.storage.UploadFile(ctx, bucket, thumbnailPath, thumbnailName)
	if err != nil {
		return result, fmt.Errorf("error uploading compressed file: %w", err)
	}

	thumbailSize, err := utils.GetFileSizeInBytes(thumbnailPath)
	if err != nil {
		return result, err
	}

	result, err = s.repository.InsertFile(ctx, model.File{
		Uri:          uri,
		ThumbnailUri: thumbnailUri,
	})

	if err != nil {
		return result, fmt.Errorf("error inserting file to database: %w", err)
	}

	log.Printf("original (%d): %s | compressed (%d): %s", sizeInBytes, uri, thumbailSize, thumbnailUri)
	return result, nil
}

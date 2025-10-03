package utils

import (
	"PattyWagon/internal/constants"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func ValidateFileExtensions(filename string, allowedExtensions []string) error {
	ext := filepath.Ext(filename)
	if slices.Contains(allowedExtensions, ext) {
		return nil
	}
	return constants.ErrInvalidFileType
}

func GetFileSizeInBytes(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, fmt.Errorf("error getting file info: %v\n", err)

	}

	fileSize := fileInfo.Size()

	return fileSize, nil
}

func ValidateFileSize(sizeInBytes int64) error {
	if sizeInBytes > constants.MaxUploadSizeInBytes {
		return fmt.Errorf("exceed max size: %v", sizeInBytes)
	}

	if sizeInBytes < constants.MinUploadSizeInBytes {
		return fmt.Errorf("less than min size: %v", sizeInBytes)
	}

	return nil
}

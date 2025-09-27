package imagecompressor

import (
	"WeenieHut/observability"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp" // For decoding
)

var (
	MaxConcurrentCompress, _ = strconv.Atoi(os.Getenv("MAX_CONCURRENT_COMPRESS"))
	CompressionQuality, _    = strconv.Atoi(os.Getenv("COMPRESSION_QUALITY"))
)

type ImageCompressor struct {
	semaphore  chan struct{}
	quality    int
	bufferPool sync.Pool
}

func New(
	maxConcurrentCompress int,
	quality int,
) *ImageCompressor {
	return &ImageCompressor{
		semaphore: make(chan struct{}, maxConcurrentCompress),
		quality:   quality,
		bufferPool: sync.Pool{New: func() any {
			return make([]byte, 0, 64*1024)
		}},
	}
}

// func (cmp *ImageCompressor) compressPNG(ctx context.Context, img image.Image) ([]byte, error) {
// 	var buf bytes.Buffer
// 	encoder := png.Encoder{
// 		CompressionLevel: png.BestCompression,
// 	}

// 	if err := encoder.Encode(&buf, img); err != nil {
// 		return nil, fmt.Errorf("error encoding image: %w", err)
// 	}

// 	return buf.Bytes(), nil
// }

func (cmp *ImageCompressor) compressJPEG(ctx context.Context, img image.Image) ([]byte, error) {
	_, span := observability.Tracer.Start(ctx, "image_compressor.compress_jpg")
	defer span.End()
	var buf bytes.Buffer

	err := jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: cmp.quality,
	})

	if err != nil {
		return nil, fmt.Errorf("error encoding image: %w", err)
	}

	return buf.Bytes(), nil
}

func (cmp *ImageCompressor) thumbnail(ctx context.Context, img image.Image, sizeInPixels int) image.Image {
	_, span := observability.Tracer.Start(ctx, "image_compressor.thumbnail")
	defer span.End()

	thumbnail := resize.Thumbnail(uint(sizeInPixels), uint(sizeInPixels), img, resize.Lanczos2)
	return thumbnail
}

func (cmp *ImageCompressor) Compress(ctx context.Context, src string) (string, error) {
	ctx, span := observability.Tracer.Start(ctx, "image_compressor.compress")
	defer span.End()

	select {
	case cmp.semaphore <- struct{}{}:
		defer func() { <-cmp.semaphore }()
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("compression queue timeout")
	}

	img, format, err := cmp.loadImage(ctx, src)
	if err != nil {
		return "", fmt.Errorf("error decoding file: %w", err)
	}

	thumbnail := cmp.thumbnail(ctx, img, 150)
	var result []byte
	switch format {
	case "jpeg":
		result, err = cmp.compressJPEG(ctx, thumbnail)
	case "png":
		result, err = cmp.compressJPEG(ctx, thumbnail)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}

	if err != nil {
		return "", fmt.Errorf("error in compressing image: %w", err)
	}

	ext := filepath.Ext(src)
	nameWithoutExt := strings.TrimSuffix(src, ext)
	resultFilename := fmt.Sprintf("%s_compressed.%s", nameWithoutExt, "jpeg")

	resultFile, err := os.Create(resultFilename)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer resultFile.Close()

	_, err = resultFile.Write(result)
	if err != nil {
		return "", fmt.Errorf("error writing data: %w", err)
	}
	return resultFilename, nil
}

func (cmp *ImageCompressor) loadImage(ctx context.Context, src string) (image.Image, string, error) {
	_, span := observability.Tracer.Start(ctx, "image_compressor.load_image")
	defer span.End()

	file, err := os.Open(src)
	if err != nil {
		return nil, "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	return image.Decode(file)
}

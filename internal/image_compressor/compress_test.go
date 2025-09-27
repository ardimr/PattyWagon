package imagecompressor

import (
	"WeenieHut/internal/utils"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	testdataDir := "testdata"
	compressor := New(5, 50)

	// t.Run("PNG", func(t *testing.T) {
	// 	input := fmt.Sprintf("%s/sample-png.png", testdataDir)
	// 	result, err := compressor.Compress(context.TODO(), input)
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, result)

	// 	originalSize, _ := utils.GetFileSizeInBytes(input)
	// 	thumbnailSize, _ := utils.GetFileSizeInBytes(result)
	// 	ratio := 100 * (float64(thumbnailSize) / float64(originalSize))
	// 	fmt.Printf("original: %d | thumbnail: %d (%.2f %%)\n", originalSize, thumbnailSize, ratio)
	// })

	t.Run("JPEG", func(t *testing.T) {
		input := fmt.Sprintf("%s/sample.jpeg", testdataDir)
		result, err := compressor.Compress(context.TODO(), input)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)

		originalSize, _ := utils.GetFileSizeInBytes(input)
		thumbnailSize, _ := utils.GetFileSizeInBytes(result)
		ratio := 100 * (float64(thumbnailSize) / float64(originalSize))
		fmt.Printf("original: %d | thumbnail: %d (%.2f %%) \n", originalSize, thumbnailSize, ratio)
	})

	// t.Run("JPG", func(t *testing.T) {
	// 	input := fmt.Sprintf("%s/sample.jpg", testdataDir)
	// 	result, err := compressor.Compress(context.TODO(), input)
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, result)
	// })
}

package location

import (
	"PattyWagon/internal/model"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber/h3-go/v4"
)

func TestService_GetAllCellIDs(t *testing.T) {
	s := &Service{}
	ctx := context.Background()

	t.Run("valid location returns 9 cells", func(t *testing.T) {
		location := model.Location{
			Lat:  -6.2088,
			Long: 106.8456,
		}

		cells, err := s.GetAllCellIDs(ctx, location)

		require.NoError(t, err)
		assert.Len(t, cells, 9, "should return 9 cells for resolutions 0-8")

		// Verify each resolution is present
		for i := 0; i <= 8; i++ {
			assert.Equal(t, i, cells[i].Resolution, "resolution should match index")
			assert.NotZero(t, cells[i].CellID, "cell ID should not be zero")
		}
	})

	t.Run("different locations return different cells", func(t *testing.T) {
		location1 := model.Location{Lat: 0.0, Long: 0.0}
		location2 := model.Location{Lat: 40.7128, Long: -74.0060}

		cells1, err1 := s.GetAllCellIDs(ctx, location1)
		cells2, err2 := s.GetAllCellIDs(ctx, location2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, cells1[8].CellID, cells2[8].CellID, "different locations should have different cell IDs")
	})

	t.Run("cell IDs increase in specificity with resolution", func(t *testing.T) {
		location := model.Location{Lat: 37.7749, Long: -122.4194}

		cells, err := s.GetAllCellIDs(ctx, location)

		require.NoError(t, err)

		// Higher resolution cells should be children of lower resolution cells
		for i := 1; i <= 8; i++ {
			parent := h3.Cell(cells[i-1].CellID)
			child := h3.Cell(cells[i].CellID)

			// Verify child is actually within parent's area
			childParent, _ := child.Parent(i - 1)
			assert.Equal(t, parent, childParent, "cell at resolution %d should be parent of resolution %d", i-1, i)
		}
	})

	t.Run("extreme coordinates", func(t *testing.T) {
		testCases := []struct {
			name     string
			location model.Location
		}{
			{"north pole", model.Location{Lat: 90.0, Long: 0.0}},
			{"south pole", model.Location{Lat: -90.0, Long: 0.0}},
			{"date line", model.Location{Lat: 0.0, Long: 180.0}},
			{"negative date line", model.Location{Lat: 0.0, Long: -180.0}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cells, err := s.GetAllCellIDs(ctx, tc.location)
				require.NoError(t, err)
				assert.Len(t, cells, 9)
			})
		}
	})
}

func TestService_FindCellIDByResolution(t *testing.T) {
	s := &Service{}
	ctx := context.Background()

	t.Run("valid resolution returns correct cell", func(t *testing.T) {
		location := model.Location{Lat: -6.2088, Long: 106.8456}
		resolution := 8

		cell, err := s.FindCellIDByResolution(ctx, location, resolution)

		require.NoError(t, err)
		assert.Equal(t, resolution, cell.Resolution)
		assert.NotZero(t, cell.CellID)
	})

	t.Run("resolution 0 returns largest cell", func(t *testing.T) {
		location := model.Location{Lat: 0.0, Long: 0.0}

		cell, err := s.FindCellIDByResolution(ctx, location, 0)

		require.NoError(t, err)
		assert.Equal(t, 0, cell.Resolution)
		assert.NotZero(t, cell.CellID)
	})

	t.Run("resolution 15 returns smallest cell", func(t *testing.T) {
		location := model.Location{Lat: 37.7749, Long: -122.4194}

		cell, err := s.FindCellIDByResolution(ctx, location, 15)

		require.NoError(t, err)
		assert.Equal(t, 15, cell.Resolution)
		assert.NotZero(t, cell.CellID)
	})

	t.Run("same location and resolution returns same cell", func(t *testing.T) {
		location := model.Location{Lat: 40.7128, Long: -74.0060}
		resolution := 10

		cell1, err1 := s.FindCellIDByResolution(ctx, location, resolution)
		cell2, err2 := s.FindCellIDByResolution(ctx, location, resolution)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, cell1.CellID, cell2.CellID)
		assert.Equal(t, cell1.Resolution, cell2.Resolution)
	})

	t.Run("different resolutions return different cells", func(t *testing.T) {
		location := model.Location{Lat: 51.5074, Long: -0.1278}

		cell5, err5 := s.FindCellIDByResolution(ctx, location, 5)
		cell10, err10 := s.FindCellIDByResolution(ctx, location, 10)

		require.NoError(t, err5)
		require.NoError(t, err10)
		assert.NotEqual(t, cell5.CellID, cell10.CellID)
		assert.Equal(t, 5, cell5.Resolution)
		assert.Equal(t, 10, cell10.Resolution)
	})
}

func TestService_FindKRingCellIDs(t *testing.T) {
	s := &Service{}
	ctx := context.Background()

	t.Run("successful k-ring generation with k=0", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := 9
		k := 0

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
		assert.Len(t, cells, 1) // k=0 should return only center cell
		assert.Equal(t, resolution, cells[0].Resolution)
	})

	t.Run("successful k-ring generation with k=1", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := 9
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
		assert.Len(t, cells, 7) // k=1 returns 7 hexagons (1 center + 6 neighbors)

		for _, cell := range cells {
			assert.Equal(t, resolution, cell.Resolution)
			assert.NotZero(t, cell.CellID)
		}
	})

	t.Run("successful k-ring generation with k=2", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := 9
		k := 2

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
		assert.Len(t, cells, 19) // k=2 returns 19 hexagons

		for _, cell := range cells {
			assert.Equal(t, resolution, cell.Resolution)
			assert.NotZero(t, cell.CellID)
		}
	})

	t.Run("different resolutions produce different cell counts", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		k := 1

		cells7, err := s.FindKRingCellIDs(ctx, location, 7, k)
		require.NoError(t, err)

		cells10, err := s.FindKRingCellIDs(ctx, location, 10, k)
		require.NoError(t, err)

		// Both should have same number of cells for same k
		assert.Len(t, cells7, 7)
		assert.Len(t, cells10, 7)

		// But different resolutions
		assert.Equal(t, 7, cells7[0].Resolution)
		assert.Equal(t, 10, cells10[0].Resolution)
	})

	t.Run("location at north pole", func(t *testing.T) {
		location := model.Location{
			Lat:  89.9,
			Long: 0,
		}
		resolution := 5
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
	})

	t.Run("location at south pole", func(t *testing.T) {
		location := model.Location{
			Lat:  -89.9,
			Long: 0,
		}
		resolution := 5
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
	})

	t.Run("location crossing dateline", func(t *testing.T) {
		location := model.Location{
			Lat:  0,
			Long: 179.9,
		}
		resolution := 6
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
		assert.Len(t, cells, 7)
	})

	t.Run("invalid resolution - too low", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := -1
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		assert.Error(t, err)
		assert.Nil(t, cells)
	})

	t.Run("invalid resolution - too high", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := 16
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		assert.Error(t, err)
		assert.Nil(t, cells)
	})

	t.Run("zero location (null island)", func(t *testing.T) {
		location := model.Location{
			Lat:  0,
			Long: 0,
		}
		resolution := 9
		k := 1

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
		assert.Len(t, cells, 7)
	})

	t.Run("verify cell IDs are unique", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := 9
		k := 2

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)

		cellIDs := make(map[int64]bool)
		for _, cell := range cells {
			assert.False(t, cellIDs[cell.CellID], "duplicate cell ID found")
			cellIDs[cell.CellID] = true
		}
	})

	t.Run("large k value", func(t *testing.T) {
		location := model.Location{
			Lat:  37.7749,
			Long: -122.4194,
		}
		resolution := 7
		k := 10

		cells, err := s.FindKRingCellIDs(ctx, location, resolution, k)

		require.NoError(t, err)
		assert.NotEmpty(t, cells)
		// k=10 should return 331 cells (3k^2 + 3k + 1)
		expectedCount := 3*k*k + 3*k + 1
		assert.Len(t, cells, expectedCount)
	})
}

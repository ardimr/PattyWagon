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

	t.Run("k=0 returns only center cell", func(t *testing.T) {
		location := model.Location{Lat: -6.2088, Long: 106.8456}

		cells, err := s.FindKRingCellIDs(ctx, location, 0)

		require.NoError(t, err)
		assert.Len(t, cells, 1, "k=0 should return only the center cell")
		assert.Equal(t, 0, cells[0].Resolution)
	})

	t.Run("k=1 returns center and immediate neighbors", func(t *testing.T) {
		location := model.Location{Lat: 0.0, Long: 0.0}

		cells, err := s.FindKRingCellIDs(ctx, location, 1)

		require.NoError(t, err)
		assert.Greater(t, len(cells), 1, "k=1 should return center plus neighbors")
		// Typically 7 cells for k=1 (1 center + 6 neighbors in hexagon)
		assert.LessOrEqual(t, len(cells), 7)
	})

	t.Run("k=2 returns larger ring", func(t *testing.T) {
		location := model.Location{Lat: 37.7749, Long: -122.4194}

		cells, err := s.FindKRingCellIDs(ctx, location, 2)

		require.NoError(t, err)
		assert.Greater(t, len(cells), 7, "k=2 should return more cells than k=1")
		// Typically 19 cells for k=2
		assert.LessOrEqual(t, len(cells), 19)
	})

	t.Run("all cells have same resolution", func(t *testing.T) {
		location := model.Location{Lat: 40.7128, Long: -74.0060}

		cells, err := s.FindKRingCellIDs(ctx, location, 1)

		require.NoError(t, err)
		require.Greater(t, len(cells), 0)

		expectedResolution := cells[0].Resolution
		for _, cell := range cells {
			assert.Equal(t, expectedResolution, cell.Resolution, "all cells should have same resolution")
		}
	})

	t.Run("larger k returns more cells", func(t *testing.T) {
		location := model.Location{Lat: 51.5074, Long: -0.1278}

		cells1, err1 := s.FindKRingCellIDs(ctx, location, 1)
		cells2, err2 := s.FindKRingCellIDs(ctx, location, 2)
		cells3, err3 := s.FindKRingCellIDs(ctx, location, 3)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)

		assert.Less(t, len(cells1), len(cells2), "k=2 should return more cells than k=1")
		assert.Less(t, len(cells2), len(cells3), "k=3 should return more cells than k=2")
	})

	t.Run("all cell IDs are unique", func(t *testing.T) {
		location := model.Location{Lat: -6.2088, Long: 106.8456}

		cells, err := s.FindKRingCellIDs(ctx, location, 2)

		require.NoError(t, err)

		cellIDMap := make(map[int64]bool)
		for _, cell := range cells {
			assert.False(t, cellIDMap[cell.CellID], "cell ID should be unique: %d", cell.CellID)
			cellIDMap[cell.CellID] = true
		}
	})

	t.Run("center cell is included in results", func(t *testing.T) {
		location := model.Location{Lat: 35.6762, Long: 139.6503}

		cells, err := s.FindKRingCellIDs(ctx, location, 2)

		require.NoError(t, err)

		// Get the expected center cell
		latLng := h3.NewLatLng(location.Lat, location.Long)
		expectedCenter, _ := h3.LatLngToCell(latLng, 0)

		found := false
		for _, cell := range cells {
			if cell.CellID == int64(expectedCenter) {
				found = true
				break
			}
		}
		assert.True(t, found, "center cell should be included in results")
	})
}

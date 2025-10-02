package location

import (
	"PattyWagon/internal/model"
	"context"

	"github.com/uber/h3-go/v4"
)

func (s *Service) GetAllCellIDs(ctx context.Context, location model.Location) ([]model.Cell, error) {
	latLng := h3.NewLatLng(location.Lat, location.Long)

	result := make([]model.Cell, 0, 9)

	for resolution := 0; resolution <= 8; resolution++ {
		cell, err := h3.LatLngToCell(latLng, resolution)
		if err != nil {
			return nil, err
		}

		result = append(result, model.Cell{
			CellID:     int64(cell),
			Resolution: resolution,
		})
	}

	return result, nil
}

func (s *Service) FindCellIDByResolution(ctx context.Context, location model.Location, resolution int) (model.Cell, error) {
	latLng := h3.NewLatLng(location.Lat, location.Long)
	cell, err := h3.LatLngToCell(latLng, resolution)
	if err != nil {
		return model.Cell{}, err
	}
	return model.Cell{
		CellID:     int64(cell),
		Resolution: cell.Resolution(),
	}, nil
}

func (s *Service) FindKRingCellIDs(ctx context.Context, location model.Location, resolution, k int) ([]model.Cell, error) {
	latLng := h3.NewLatLng(location.Lat, location.Long)
	centerCell, err := h3.LatLngToCell(latLng, resolution)
	if err != nil {
		return nil, err
	}

	cells, err := h3.GridDisk(centerCell, k)
	if err != nil {
		return nil, err
	}

	result := make([]model.Cell, 0, len(cells))
	for _, cell := range cells {
		result = append(result, model.Cell{
			CellID:     int64(cell),
			Resolution: cell.Resolution(),
		})
	}

	return result, nil
}

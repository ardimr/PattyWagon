package utils

import (
	"PattyWagon/internal/constants"
	"strconv"
	"strings"
)

func ValidateAndExtractCoordinate(raw string) (float64, float64, error) {
	temp := strings.Split(strings.ReplaceAll(raw, " ", ""), ",")
	if len(temp) < 2 {
		return 0, 0, constants.ErrInvalidCoordinate
	}

	lat, err := strconv.ParseFloat(temp[0], 64)
	if err != nil {
		return 0, 0, constants.ErrInvalidCoordinate
	}

	long, err := strconv.ParseFloat(temp[1], 64)
	if err != nil {
		return 0, 0, constants.ErrInvalidCoordinate
	}

	return lat, long, nil
}

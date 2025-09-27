package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAndExtractCoordinate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedLat float64
		expectedLng float64
		expectError bool
	}{
		{
			name:        "valid coordinates",
			input:       "40.7128,-74.0060",
			expectedLat: 40.7128,
			expectedLng: -74.0060,
			expectError: false,
		},
		{
			name:        "valid coordinates with spaces",
			input:       "40.7128, -74.0060",
			expectedLat: 40.7128,
			expectedLng: -74.0060,
			expectError: false,
		},
		{
			name:        "zero coordinates",
			input:       "0,0",
			expectedLat: 0,
			expectedLng: 0,
			expectError: false,
		},
		{
			name:        "positive coordinates",
			input:       "51.5074,0.1278",
			expectedLat: 51.5074,
			expectedLng: 0.1278,
			expectError: false,
		},
		{
			name:        "extreme valid coordinates",
			input:       "90,-180",
			expectedLat: 90,
			expectedLng: -180,
			expectError: false,
		},
		{
			name:        "missing comma",
			input:       "40.7128-74.0060",
			expectedLat: 0,
			expectedLng: 0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectedLat: 0,
			expectedLng: 0,
			expectError: true,
		},
		{
			name:        "only latitude",
			input:       "40.7128",
			expectedLat: 0,
			expectedLng: 0,
			expectError: true,
		},
		{
			name:        "invalid latitude",
			input:       "invalid,-74.0060",
			expectedLat: 0,
			expectedLng: 0,
			expectError: true,
		},
		{
			name:        "invalid longitude",
			input:       "40.7128,invalid",
			expectedLat: 0,
			expectedLng: 0,
			expectError: true,
		},
		{
			name:        "too many parts",
			input:       "40.7128,-74.0060,extra",
			expectedLat: 40.7128,
			expectedLng: -74.0060,
			expectError: false,
		},
		{
			name:        "latitude out of range",
			input:       "91.0,-74.0060",
			expectedLat: 91.0,
			expectedLng: -74.0060,
			expectError: false, // function doesn't validate ranges
		},
		{
			name:        "longitude out of range",
			input:       "40.7128,-181.0",
			expectedLat: 40.7128,
			expectedLng: -181.0,
			expectError: false, // function doesn't validate ranges
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng, err := ValidateAndExtractCoordinate(tt.input)

			if tt.expectError {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.expectedLat, lat)
			assert.Equal(t, tt.expectedLng, lng)

		})
	}
}

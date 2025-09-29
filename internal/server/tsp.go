import (
	"math"
)

// {"A",1,2}
// {"B",2,5}
// {"C",4,1}
// {"X",7,8}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	deltaLat := radians(lat2 - lat1)
	deltaLon := radians(lon2 - lon1)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + (math.Cos(radians(lat1)) * math.Cos(radians(lat2)) * math.Sin(deltaLon/2) * math.Sin(deltaLon/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return d
}

func radians(deg float64) float64 {
	return deg * math.Pi / 180
}


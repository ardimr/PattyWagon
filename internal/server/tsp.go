import (
	"math"
)

type Location struct {
	Name string
	Lat  float64
	Lon  float64
}

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

func toDistanceGraph() map[string]map[string]float64 {

	// assume the location is from variabel input
	var locations = []Location{
		{"Merchant A", 40.7128, -74.0060},
		{"Merchant B", 37.1234, -122.6543},
		{"Merchant C", -12.8756, 45.1234},
		{"User D", 51.5076, -0.1227},
	}
	n := len(locations)
	graph := make(map[string]map[string]float64)
	for i := 0; i < n; i++ {
		dist := make(map[string]float64)
		for j := 0; j < n; j++ {
			if i != j {
				dist[locations[j].Name] = haversine(locations[i].Lat, locations[i].Lon, locations[j].Lat, locations[j].Lon)
			}
		}
		graph[locations[i].Name] = dist
	}
	return graph
}

// generate permutations of merchant
func permutations(arr []string) [][]string {
	var helper func([]string, int)
	res := [][]string{}

	helper = func(a []string, i int) {
		if i == len(a)-1 {
			tmp := make([]string, len(a))
			copy(tmp, a)
			res = append(res, tmp)
			return
		}
		for j := i; j < len(a); j++ {
			a[i], a[j] = a[j], a[i]
			helper(a, i+1)
			a[i], a[j] = a[j], a[i]
		}
	}

	helper(arr, 0)
	return res
}


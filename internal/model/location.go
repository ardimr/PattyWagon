package model

type Cell struct {
	CellID     int64
	Resolution int
}

type Location struct {
	Lat             float64
	Long            float64
	IsStartingPoint bool
}

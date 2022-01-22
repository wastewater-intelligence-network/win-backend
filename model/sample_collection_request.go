package model

type Location struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

type SampleCollectionRequest struct {
	ContainerId string   `json:"container_id"`
	Location    Location `json:"location"`
}

package model

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type SampleCollectionRequest struct {
	ContainerId string   `json:"container_id"`
	Location    Location `json:"location"`
}

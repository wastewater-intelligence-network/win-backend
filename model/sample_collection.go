package model

import "time"

type SampleStatus string

const (
	SampleStatusCollected     SampleStatus = "sample_collected"
	SampleStatusInTransit     SampleStatus = "sample_in_transit"
	SampleStatusReceivedAtLab SampleStatus = "sample_received_in_lab"
	SampleStatusTestInProgess SampleStatus = "sample_test_in_progress"
	SampleStatusResultOut     SampleStatus = "sample_result_out"
)

type Location struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

type StatusLog struct {
	Timestamp time.Time
	Status    SampleStatus
	Message   string
}

type Sample struct {
	SampleTakenOn            time.Time       `json:"sampleTakenOn,omitempty" bson:"sampleTakenOn"`
	ContainerId              string          `json:"containerId,omitempty" bson:"containerId"`
	SampleCollectionLocation CollectionPoint `json:"sampleCollectionLocation,omitempty" bson:"sampleCollectionLocation"`
	Status                   SampleStatus    `json:"status,omitempty" bson:"status"`
	StatusLogList            []StatusLog     `json:"statusLog,omitempty" bson:"statusLog"`
}

type SamplingRequest struct {
	ContainerId string   `json:"containerId,omitempty" bson:"containerId"`
	Location    Location `json:"location,omitempty" bson:"location"`
	PointId     string   `json:"pointId,omitempty" bson:"pointId"`
}

type SamplingStatusRequest struct {
}

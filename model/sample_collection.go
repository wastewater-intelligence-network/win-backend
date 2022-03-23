package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SampleStatus string

const (
	SampleStatusCollected     SampleStatus = "sample_collected"
	SampleStatusInTransit     SampleStatus = "sample_in_transit"
	SampleStatusReceivedAtLab SampleStatus = "sample_received_in_lab"
	SampleStatusTestInProgess SampleStatus = "sample_test_in_progress"
	SampleStatusResultOut     SampleStatus = "sample_result_out"
)

var StatusSequence = []SampleStatus{
	SampleStatusCollected,
	SampleStatusInTransit,
	SampleStatusReceivedAtLab,
	SampleStatusTestInProgess,
	SampleStatusResultOut,
}

type Location struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

type StatusLog struct {
	Timestamp time.Time    `json:"timestamp,omitempty" bson:"timestamp"`
	Status    SampleStatus `json:"status,omitempty" bson:"status"`
	Message   string       `json:"message,omitempty" bson:"message"`
	Error     string       `json:"error,omitempty" bson:"error"`
}

type Sample struct {
	SampleId                 primitive.ObjectID     `json:"sampleId,omitempty" bson:"_id,omitempty"`
	SampleTakenOn            time.Time              `json:"sampleTakenOn,omitempty" bson:"sampleTakenOn"`
	ContainerId              string                 `json:"containerId,omitempty" bson:"containerId"`
	SampleCollectionLocation CollectionPoint        `json:"sampleCollectionLocation,omitempty" bson:"sampleCollectionLocation"`
	Status                   SampleStatus           `json:"status,omitempty" bson:"status"`
	StatusLogList            []StatusLog            `json:"statusLog,omitempty" bson:"statusLog"`
	AdditionalData           map[string]interface{} `json:"additionalData,omitempty" bson:"additionalData"`
}

type SamplingRequest struct {
	ContainerId    string                 `json:"containerId,omitempty" bson:"containerId"`
	Location       Location               `json:"location,omitempty" bson:"location"`
	PointId        string                 `json:"pointId,omitempty" bson:"pointId"`
	AdditionalData map[string]interface{} `json:"additionalData,omitempty" bson:"additionalData"`
}

type SamplingStatusRequest struct {
	ContainerId string `json:"containerId,omitempty" bson:"containerId"`
	StatusPatch string `json:"statusPatch,omitempty" bson:"statusPatch,omitempty"`
}

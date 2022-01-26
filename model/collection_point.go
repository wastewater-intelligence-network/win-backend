package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CollectionPoint struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PointId  string             `json:"pointId,omitempty" bson:"pointId,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Location Location           `json:"location,omitempty" bson:"location,omitempty"`
	Type     string             `json:"type,omitempty" bson:"type,omitempty"`
}

package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CollectionPoint struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Location Location           `json:"location,omitempty"`
	Type     string             `json:"type,omitempty"`
}

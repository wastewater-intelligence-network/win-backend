package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CollectionPointSchedule struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Date              string             `bson:"date,omitempty"`
	Time              string             `bson:"time,omitempty"`
	AssignedToID      primitive.ObjectID `bson:"assignedToId,omitempty"`
	CollectionPointID primitive.ObjectID `bson:"collectionPointId,omitempty"`
}

type CollectionSchedule struct {
	CollectionPointSchedule []CollectionPointSchedule `bson:"collectionPointSchedule,omitempty"`
}

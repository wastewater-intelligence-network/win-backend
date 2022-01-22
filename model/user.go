package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Type  string             `bson:"type,omitempty"`
	Email string             `bson:"email,omitempty"`
	Hash  string             `bson:"hash,omitempty"`
}

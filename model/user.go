package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Type     string             `bson:"type,omitempty"`
	Username string             `bson:"username,omitempty"`
	Hash     string             `bson:"hash,omitempty"`
	Password string             `bson:"password,omitempty"`
	Roles    []string           `bson:"roles,omitempty"`
}

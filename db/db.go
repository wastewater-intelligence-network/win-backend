package db

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_URI = "mongodb://127.0.0.1:27017"
var DB_NAME = "win"

type DBConnection struct {
	client   *mongo.Client
	Database *mongo.Database
	ctx      context.Context
}

func NewDBConnection() (*DBConnection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(DB_URI))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &DBConnection{
		client:   client,
		Database: client.Database(DB_NAME),
		ctx:      ctx,
	}, nil
}

func (conn *DBConnection) Insert(collection string, doc interface{}) error {
	_, err := conn.Database.Collection(collection).InsertOne(conn.ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (conn *DBConnection) Find(collection string, filter interface{}) (*mongo.Cursor, error) {
	cursor, err := conn.Database.Collection(collection).Find(conn.ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (conn *DBConnection) FindOne(collection string, filter interface{}) *mongo.SingleResult {
	res := conn.Database.Collection(collection).FindOne(conn.ctx, filter)
	return res
}

func (conn *DBConnection) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	res, err := conn.Database.Collection(collection).UpdateOne(conn.ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (conn *DBConnection) DeleteCollection(collection string) error {
	_, err := conn.Database.Collection(collection).DeleteMany(conn.ctx, gin.H{})
	if err != nil {
		return err
	}
	return nil
}

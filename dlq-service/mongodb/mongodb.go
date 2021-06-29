package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is a collection of behavior of mongodb client.
type Client interface {
	Connect(ctx context.Context) (err error)
	Database(name string, opts ...*options.DatabaseOptions) (db Database)
	Disconnect(ctx context.Context) (err error)
}

// Database is a collection of behavior of mongodb database.
type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) (col Collection)
}

// Collection is a collection of behavior of mongodb client.
type Collection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (result SingleResult)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cursor Cursor, err error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (result *mongo.InsertOneResult, err error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (result *mongo.InsertManyResult, err error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (counted int64, err error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error)
}

// SingleResult is a collectioin of function of mongodb single result.
type SingleResult interface {
	Decode(v interface{}) error
	Err() error
}

// Cursor is a collection of function of mongodb cursor.
type Cursor interface {
	ID() int64
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}

package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseAdapter is concrete struct of mongodb database adapter.
type DatabaseAdapter struct {
	db *mongo.Database
}

// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
func (db *DatabaseAdapter) Collection(name string, opts ...*options.CollectionOptions) (col Collection) {
	mongoCollection := db.db.Collection(name, opts...)
	col = &CollectionAdapter{mongoCollection}
	return
}

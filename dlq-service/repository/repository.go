package repository

import (
	"context"

	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/entity"
	"github.com/sangianpatrick/go-kafka-dlq-demo/dlq-service/mongodb"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DLQRepository interface {
	FindMany(ctx context.Context, limit, skip int64) (bunchOfMessage []entity.Message, err error)
	FindByID(ctx context.Context, ID string) (message entity.Message, err error)
	DeleteByID(ctx context.Context, ID string) (err error)
	CountDocuments(ctx context.Context) (counted int64, err error)
}

type dlqRepository struct {
	logger     *logrus.Logger
	collection string
	db         mongodb.Database
}

func NewDLQRepository(logger *logrus.Logger, db mongodb.Database) DLQRepository {
	collection := "dlq-message"
	return &dlqRepository{
		logger:     logger,
		collection: collection,
		db:         db,
	}
}

func (r *dlqRepository) CountDocuments(ctx context.Context) (counted int64, err error) {
	filter := bson.M{}
	options := options.Count()

	counted, err = r.db.Collection(r.collection).CountDocuments(ctx, filter, options)

	return
}

func (r *dlqRepository) FindMany(ctx context.Context, limit, skip int64) (bunchOfMessage []entity.Message, err error) {

	filter := bson.M{}
	options := options.Find().SetLimit(limit).SetSkip(skip)

	cursor, err := r.db.Collection(r.collection).Find(ctx, filter, options)
	if err != nil {
		r.logger.Error(err)
		return
	}

	for cursor.Next(ctx) {
		var message entity.Message

		if err = cursor.Decode(&message); err != nil {
			r.logger.Error(err)
			return
		}

		bunchOfMessage = append(bunchOfMessage, message)
	}

	return
}

func (r *dlqRepository) FindByID(ctx context.Context, ID string) (message entity.Message, err error) {

	filter := bson.M{
		"id": ID,
	}
	options := options.FindOne()

	result := r.db.Collection(r.collection).FindOne(ctx, filter, options)

	if err = result.Decode(&message); err != nil {
		r.logger.Error(err)
		return
	}

	return
}

func (r *dlqRepository) DeleteByID(ctx context.Context, ID string) (err error) {

	filter := bson.M{
		"id": ID,
	}
	options := options.Delete()

	_, err = r.db.Collection(r.collection).DeleteOne(ctx, filter, options)

	return
}

package bucket

import (
	"context"
	"node/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (b Bucket) Insert(ctx context.Context, collection string, id string, value Document) (time.Time, error) {

	coll := b.client.Database(config.GetDatabaseName()).Collection(collection)
	insertTime := time.Now()

	doc := bson.M{
		"_id":       id,
		"content":   value.Content,
		"timestamp": insertTime,
	}

	_, err := coll.InsertOne(ctx, doc)

	if mongo.IsDuplicateKeyError(err) {
		// if the key exists already, update the record
		_, err := coll.ReplaceOne(ctx, bson.M{"_id": id}, doc)
		if err != nil {
			return time.Time{}, err
		}

		b.logger.Debug("updated", zap.Any("updated_id", id))
	} else {
		// if this wasn't duplicate key error, return
		if err != nil {
			return time.Time{}, err
		}

		// another case is that err == nil -> inserted successfully
		b.logger.Debug("inserted", zap.Any("inserted_id", id))
	}

	return insertTime, nil
}

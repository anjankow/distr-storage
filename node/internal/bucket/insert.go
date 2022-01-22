package bucket

import (
	"context"
	"node/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return time.Time{}, err
	}

	b.logger.Debug("inserted", zap.Any("inserted_id", result.InsertedID))

	return insertTime, nil
}

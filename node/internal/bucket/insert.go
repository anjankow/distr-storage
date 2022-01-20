package bucket

import (
	"context"
	"node/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (b Bucket) Insert(ctx context.Context, collection string, id string, value Document) error {

	coll := b.client.Database(config.GetDatabaseName()).Collection(collection)

	doc := bson.M{
		"_id":     id,
		"content": value.Content,
	}

	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	b.logger.Debug("inserted", zap.Any("inserted_id", result.InsertedID))

	return nil
}

package bucket

import (
	"context"
	"encoding/json"
	"node/internal/config"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

func (b Bucket) Insert(ctx context.Context, collection string, key string, value json.RawMessage) error {

	coll := b.client.Database(config.GetDatabaseName()).Collection(collection)

	doc := bson.D{{"_id", key}, {"value", value}}
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	b.logger.Debug("inserted", zap.String("key", key), zap.Any("inserted_id", result.InsertedID))

	return nil
}

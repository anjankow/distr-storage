package bucket

import (
	"context"
	"errors"
	"fmt"
	"node/internal/config"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

func (b Bucket) Delete(ctx context.Context, collection string, id string) error {

	coll := b.client.Database(config.GetDatabaseName()).Collection(collection)

	filters := bson.M{
		"_id": id,
	}

	result, err := coll.DeleteOne(ctx, filters)
	if err != nil {
		return errors.New("failed to delete: " + err.Error())
	}

	if result.DeletedCount > 1 {
		return fmt.Errorf("deleted unexpected number of elements: %d", result.DeletedCount)
	}

	if result.DeletedCount == 0 {
		b.logger.Debug("no documents removed", zap.String("id", id))
	}

	return nil
}

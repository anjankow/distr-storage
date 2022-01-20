package bucket

import (
	"context"
	"node/internal/config"

	"gopkg.in/mgo.v2/bson"
)

func (b Bucket) Get(ctx context.Context, collection string, id string) (Document, error) {

	coll := b.client.Database(config.GetDatabaseName()).Collection(collection)

	filters := bson.M{
		"_id": id,
	}

	var doc Document
	err := coll.FindOne(ctx, filters).Decode(&doc)

	return doc, err
}

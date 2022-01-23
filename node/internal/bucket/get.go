package bucket

import (
	"context"
	"errors"
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

func (b Bucket) GetAll(ctx context.Context, collection string) ([]Document, error) {

	coll := b.client.Database(config.GetDatabaseName()).Collection(collection)

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("bucket: failed to find all documents: " + err.Error())
	}
	var docs []Document
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.New("bucket: failed to get all docs from the cursor: " + err.Error())
	}

	return docs, nil
}

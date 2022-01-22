package app

import (
	"context"
	"encoding/json"
	"node/internal/bucket"
	"node/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	bucket bucket.Bucket
}

func NewApp(l *zap.Logger) (app App, err error) {

	bucket, err := bucket.NewConnection(l, config.GetDbConnectionURI())
	if err != nil {
		l.Fatal("db connection not established: " + err.Error())
	}

	app = App{
		Logger: l,
		bucket: bucket,
	}
	return
}

func (a App) Insert(ctx context.Context, collection string, id string, content json.RawMessage) (time.Time, error) {
	// save in the dedicated bucket
	doc := bucket.Document{
		Content: content,
	}
	return a.bucket.Insert(ctx, collection, id, doc)
}

func (a App) Get(ctx context.Context, collection string, id string) (json.RawMessage, error) {
	result, err := a.bucket.Get(ctx, collection, id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.Logger.Warn("no documents found", zap.String("id", id))
			return nil, nil
		}

		return nil, err
	}

	a.Logger.Debug("found", zap.String("id", id))

	return result.Content, nil
}

func (a App) Delete(ctx context.Context, collection string, id string) error {

	return a.bucket.Delete(ctx, collection, id)
}

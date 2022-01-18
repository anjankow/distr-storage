package app

import (
	"context"
	"encoding/json"
	"node/internal/bucket"
	"node/internal/config"

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

func (a App) Insert(ctx context.Context, collection string, key string, value json.RawMessage) error {
	// save in the dedicated bucket
	return a.bucket.Insert(ctx, collection, key, value)
}

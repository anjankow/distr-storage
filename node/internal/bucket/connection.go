package bucket

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Bucket struct {
	// connection closer function
	Disconnect func()

	client *mongo.Client
	logger *zap.Logger
}

func NewConnection(logger *zap.Logger, uri string) (Bucket, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("db connection failed", zap.String("uri", uri))
		return Bucket{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return Bucket{}, err
	}

	closer := func() {
		if err = client.Disconnect(context.Background()); err != nil {
			logger.Error("failed to disconnect the DB: " + err.Error())
		}
	}

	return Bucket{
		Disconnect: closer,
		client:     client,
		logger:     logger,
	}, nil

}

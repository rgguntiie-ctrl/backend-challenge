package infrastructure

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.uber.org/zap"
)

func NewMongoClient(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		zap.L().Fatal("failed to connect to mongo", zap.Error(err))
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		zap.L().Fatal("failed to ping mongo", zap.Error(err))
	}
	return client
}

func MongoDisconnect(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		zap.L().Fatal("failed to disconnect mongo client", zap.Error(err))
	}
}

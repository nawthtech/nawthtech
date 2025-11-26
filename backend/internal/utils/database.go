package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

func InitDatabase(databaseURL string) error {
	if databaseURL == "" {
		return nil // لا توجد قاعدة بيانات، ربما بيئة تطوير
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(databaseURL)
	
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// التحقق من الاتصال
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	db = client.Database("nawthtech")
	return nil
}

func GetDB() *mongo.Database {
	return db
}

func GetCollection(name string) *mongo.Collection {
	if db == nil {
		return nil
	}
	return db.Collection(name)
}

func CloseDatabase() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client.Disconnect(ctx)
	}
}

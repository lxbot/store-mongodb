package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type M = map[string]interface{}

var ch *chan M
var collection *mongo.Collection

func Boot(c *chan M) {
	ch = c
	uri := os.Getenv("LXBOT_MONGODB_URI")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	collection = client.Database("lxbot").Collection("store")
	// TODO: create index
}

func Set(key string, value interface{}) {
	filter := bson.M{"key": key}
	update := bson.M{
		"key": key,
		"value": value,
	}
	upsert := true
	option := options.UpdateOptions{Upsert: &upsert}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, _ = collection.UpdateOne(ctx, filter, update, &option)
}

func Get(key string) interface{} {
	filter := bson.M{"key": key}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := collection.FindOne(ctx, filter)
	m := map[string]interface{}{}
	_ = result.Decode(&m)

	return m["value"]
}
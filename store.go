package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
		log.Fatalln("store-mongodb mongo connect error:", err)
	}
	collection = client.Database("lxbot").Collection("store")

	t := true
	_, _ = collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.M{"key": 1}, Options: &options.IndexOptions{Background: &t, Unique: &t}})
}

func Set(key string, value interface{}) {
	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"key": key,
			"value": value,
		},
	}
	upsert := true
	option := options.UpdateOptions{Upsert: &upsert}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := collection.UpdateOne(ctx, filter, update, &option); err != nil {
		fmt.Printf("%v", err)
	}
}

func Get(key string) interface{} {
	filter := bson.M{"key": key}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := collection.FindOne(ctx, filter)
	m := M{}
	_ = result.Decode(&m)

	return m["value"]
}

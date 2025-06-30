package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertFPintoDB(fingerPrint map[uint32]information) error {
	// Replace with your MongoDB URI
	uri := "mongodb+srv://shaizahlabique:OSVME1sLPfwH4gRl@cluster0.d9fzhp3.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	// Set connection options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Ping to test connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB!")

	// Example: get a collection
	collection := client.Database("shazam").Collection("finger_prints")

	// Example: insert a document

	for hash, info := range fingerPrint {
		filter := bson.M{"_id": hash}
		update := bson.M{
			"$push": bson.M{
				"couples": bson.M{
					"anchorTimeMs": info.anchor_time,
					"songID":       info.songID,
				},
			},
		}
		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	fmt.Println("All fingerprints inserted.")

	return nil

}

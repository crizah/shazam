package db

import (
	"errors"
	"shazam/shazam"

	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func PutintoDB(fingerPrint map[uint32]shazam.Information) error {
	uri := os.Getenv("MONGODB_URI")

	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {

		return errors.New("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)

	collection := client.Database("shazam").Collection("finger_prints") // need to change this but this will do for now

	for hash, in := range fingerPrint {

		filter := bson.M{"_id": hash}
		update := bson.M{
			"$push": bson.M{
				"couples": bson.M{
					"anchorTimeMs": in.Anchor_time,
					"songID":       in.SongID,
				},
			},
		}

		_, err := collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))

		if err != nil {
			return err
		}

	}

	fmt.Println("All fingerprints inserted.")

	return nil

}

func SearchDB(sampleHash uint32, sampleTime uint32) (uint32, uint32) {
	// returns matched hash and also matched sonngID

}

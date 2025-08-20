package db

import (
	"errors"
	// "shazam/shazam"

	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Information struct {
	Anchor_time uint32
	SongID      uint32
}

func PutintoDB(fingerPrint map[uint32]Information) error {
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

type Matched struct {
	SampleTime  uint32
	MatchedTime uint32
	DBsongId    uint32
}

func SearchDB(sampleHash uint32, sampleTime uint32) ([]Matched, error) {
	// returns matched hash and also matched sonngID
	var matches []Matched

	uri := os.Getenv("MONGODB_URI")

	// docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {

		return matches, errors.New("MONGODB_URI empty")
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))

	if err != nil {
		return matches, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)

	collection := client.Database("shazam").Collection("finger_prints")

	filter := bson.M{"_id": sampleHash}

	cursor, err := collection.Find(ctx, filter, options.Find())

	if err != nil {
		return matches, err
	}

	var found []Information

	if err = cursor.All(ctx, &found); err != nil {
		return matches, err
	}

	for _, f := range found {
		m := Matched{SampleTime: sampleTime, MatchedTime: f.Anchor_time, DBsongId: f.SongID}
		matches = append(matches, m)

	}

	if len(matches) == 0 {
		return matches, errors.New("no matches found")
	}

	return matches, nil

}

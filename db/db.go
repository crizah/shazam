package db

import (
	"errors"
	"shazam/structs"

	// "shazam/shazam/FingerPrints"

	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func PutintoDB(fingerPrint map[uint32]structs.Information) error {
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
	MatchedHash uint32
	SampleTime  uint32
	MatchedTime uint32
	DBsongId    uint32
}

func SearchDB(sampleFP map[uint32]structs.Information) (map[uint32][]Matched, error) {

	uri := os.Getenv("MONGODB_URI")

	// docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {

		return nil, errors.New("MONGODB_URI empty")
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)

	collection := client.Database("shazam").Collection("finger_prints")

	// var intermediate [](map[uint32][]Matched)
	var bins map[uint32][]Matched

	for h, info := range sampleFP {
		// per hash in the sample dingerPrint

		filter := bson.M{"_id": h}

		cursor, err := collection.Find(ctx, filter, options.Find())

		var found []structs.Information

		if err = cursor.All(ctx, &found); err != nil {
			return nil, err
		}

		// var arr []Matched
		// var bin map[uint32]Matched

		for _, f := range found {
			m := Matched{SampleTime: info.Anchor_time, MatchedTime: f.Anchor_time, DBsongId: f.SongID, MatchedHash: h}
			// need to store a map per song with id as h and value as matched

			city, ok := bins[f.SongID]
			if ok {
				bins[f.SongID] = append(city, m)

			} else {
				arr := []Matched{m}
				bins[f.SongID] = arr

			}

		}

	}

	return bins, nil

	// var bins map[uint32][]Matched

	// need to break it into maps per song

}

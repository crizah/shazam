package db

import (
	"shazam/structs"
	"time"

	// "shazam/shazam/FingerPrints"

	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

func NewMongoClient(uri string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoClient{client: client}, nil
}

func (db *MongoClient) Close() error {
	if db.client != nil {
		return db.client.Disconnect(context.Background())
	}
	return nil
}

func (db *MongoClient) PutSongIds(songId uint32, artist string, songName string) error {

	collection := db.client.Database("shazamm").Collection("test_songs") // need to change this but this will do for now

	// key as songId and name, artist

	filter := bson.M{"_id": songId}
	update := bson.M{
		"$push": bson.M{
			"couples": bson.M{
				"songName":   songName,
				"songArtist": artist,
			},
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update, options.UpdateOne().SetUpsert(true))

	if err != nil {
		return err
	}

	fmt.Println("song registered")
	return nil

}

func (db *MongoClient) PutintoDB(fingerPrint structs.OMap) error {

	collection := db.client.Database("shazamm").Collection("test_fps") // need to change this but this will do for now

	for hash, in := range fingerPrint.Map {

		filter := bson.M{"_id": hash}
		update := bson.M{
			"$push": bson.M{
				"couples": bson.M{
					"anchorTimeMs": in.Anchor_time,
					"songID":       in.SongID,
				},
			},
		}

		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		_, err := collection.UpdateOne(context.Background(), filter, update, options.UpdateOne().SetUpsert(true))

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

func (db *MongoClient) SearchDB(sampleFP structs.OMap) (map[uint32][]Matched, error) {

	collection := db.client.Database("shazamm").Collection("test_fps")

	bins := make(map[uint32][]Matched)

	for h, info := range sampleFP.Map {
		// per hash in the sample dingerPrint

		filter := bson.M{"_id": h}

		cursor, err := collection.Find(context.Background(), filter, options.Find())
		if err != nil {
			return nil, err
		}

		var found []structs.Information

		if err = cursor.All(context.Background(), &found); err != nil {
			return nil, err
		}

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

}

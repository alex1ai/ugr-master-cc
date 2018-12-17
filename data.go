package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	log "github.com/sirupsen/logrus"
	"time"
)

type DB struct {
	client *mongo.Client
}

func (db *DB) connect(ip string, port int) (err error) {
	log.Infof("Connecting to Mongo Database, make sure it is running on %s:%d", ip, port)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.client, err = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%d", ip, port))

	if err != nil {
		return err
	}

	// Test reaching the DB
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = db.client.Ping(ctx, readpref.Primary())
	return
}

func (db DB) query(query map[string]interface{}) ([] Content, error) {

	collection := db.client.Database(Database).Collection(Collection)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, query)

	if err != nil {
		return nil, err
	}

	response := make([]Content, 1)
	for cur.Next(ctx) {
		var result Content
		err = cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		response = append(response, result)
	}

	return response, nil

}

func (db DB) update(filter interface{}, replacement interface{}) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coll := db.client.Database(Database).Collection(Collection)

	opts := options.ReplaceOptions{}
	// If lang and id not there yet, this represents the same as a PUT request, i.e. creating a new Document
	opts.SetUpsert(true)
	ex, err := coll.ReplaceOne(ctx, filter, replacement, &opts)

	return ex != nil, err
}

func (db DB) delete(instanceQuery interface{}) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coll := db.client.Database(Database).Collection(Collection)

	del, err := coll.DeleteOne(ctx, instanceQuery)
	log.Debug(del.DeletedCount)

	return del, err
}

func (db DB) close() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return db.client.Disconnect(ctx)
}

func (db DB) drop() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.client.Database(Database).Drop(ctx)
	log.Debug("Dropped database")
}

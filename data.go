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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	db.client, err = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%d", ip, port))

	if err != nil {
		return err
	}

	// Test reaching the DB
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = db.client.Ping(ctx, readpref.Primary())
	return
}

func (db DB) populate(n int) error {

	collection := db.client.Database(Database).Collection(Collection)

	content := make([]interface{}, n)
	for i := 0; i < n; i++ {
		content[i] = createDummyContent(i + 1)
		time.Sleep(1 * time.Second)
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := collection.InsertMany(ctx, content)

	return err
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

func (db DB) reset() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := db.client.Database(Database).Collection(Collection).DeleteMany(ctx, nil)
	if err != nil {
		log.Error(err)
	}
	log.Debug("Deleted everything in database")
}

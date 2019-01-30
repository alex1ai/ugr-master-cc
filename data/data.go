package data

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var (
	Database          = "info"
	ContentCollection = "content"
)

const timeOut = 200 * time.Millisecond

type DB struct {
	Client *mongo.Client
}

func (db *DB) Connect(ip string, port int) (err error) {
	log.Infof("Connecting to Mongo Database, make sure it is running on %s:%d", ip, port)
	ctx, cancel := context.WithTimeout(context.Background(), timeOut*10)
	defer cancel()
	db.Client, err = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%d", ip, port))

	if err != nil {
		return err
	}

	// Test reaching the DB
	return db.Client.Ping(ctx, readpref.Primary())
}

func (db DB) Populate(n int) error {

	collection := db.Client.Database(Database).Collection(ContentCollection)

	content := make([]interface{}, n)
	for i := 0; i < n; i++ {
		content[i] = createDummyContent(i + 1)
	}

	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	_, err := collection.InsertMany(ctx, content)

	return err
}

func (db DB) Query(query map[string]interface{}) ([]Content, error) {

	collection := db.Client.Database(Database).Collection(ContentCollection)

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	cur, err := collection.Find(ctx, query)

	if err != nil {
		return nil, err
	}

	response := make([]Content, 0)
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

func getNextId(db *DB) (uint, error) {
	const startID = 1
	coll := db.Client.Database(Database).Collection(ContentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), timeOut*2)
	defer cancel()
	opts := options.FindOptions{}
	opts.Sort = bson.M{"id": -1}
	var limit int64 = 1
	opts.Limit = &limit
	cur, err := coll.Find(ctx, bson.M{}, &opts)
	if err != nil {
		return startID, err
	}
	var result Content
	for cur.Next(ctx) {
		if err = cur.Decode(&result); err != nil {
			return startID, err
		}
	}
	return result.Id + 1, nil

}

func (db DB) Add(newInstance Content) (bool, error) {
	coll := db.Client.Database(Database).Collection(ContentCollection)

	query := map[string]interface{}{
		"id": newInstance.Id,
	}
	res, err := db.Query(query)
	// Raise to next db max id if it was not set or seems to be set randomly
	if newInstance.Id == 0 || (err == nil && len(res) == 0) {
		newInstance.Id, err = getNextId(&db)
		if err != nil {
			return false, err
		}
	}
	var exists bool
	for _, c := range res{
		if c.Id == newInstance.Id && c.Language == newInstance.Language{
			exists = true
			break
		}
	}
	if exists {
		return false, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	newInstance.SetCreatedAtNow()
	ex, err := coll.InsertOne(ctx, newInstance)

	return ex != nil, err
}

func (db DB) Update(filter interface{}, replacement interface{}) (bool, error) {
	coll := db.Client.Database(Database).Collection(ContentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	opts := options.ReplaceOptions{}
	// If lang and id not there yet, this represents the same as a PUT request, i.e. creating a new Document
	opts.SetUpsert(true)
	ex, err := coll.ReplaceOne(ctx, filter, replacement, &opts)

	return ex != nil, err
}

func (db DB) Delete(instanceQuery interface{}) (*mongo.DeleteResult, error) {
	coll := db.Client.Database(Database).Collection(ContentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	del, err := coll.DeleteOne(ctx, instanceQuery)
	log.Debug(del.DeletedCount)

	return del, err
}

func (db DB) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	return db.Client.Disconnect(ctx)
}

func (db DB) Reset() {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	err := db.Client.Database(Database).Drop(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Debug("Deleted everything in database")
}

func createDummyContent(id int) Content {
	langs := []string{"de", "en", "es", "ar"}
	lang := langs[rand.Intn(len(langs))]
	created := time.Now()
	return Content{"test 1", "test1 answer", uint(id), lang, "work", created}
}

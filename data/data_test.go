package data

import (
	"os"
	"strconv"
	"testing"
)

var db *DB

func setupDB(t *testing.T) {
	if db == nil {
		Database = "testing"
		data := DB{}

		mPort := getEnv("MONGO_PORT", "27017")
		mIP := getEnv("MONGO_IP", "localhost")

		portI, err := strconv.Atoi(mPort)
		err = data.Connect(mIP, portI)
		if err != nil {
			t.Error(err)
		}
		if err != nil {
			t.Fatal(err)
		}
		db = &data
	}
}


func TestDB_Add(t *testing.T) {
	setupDB(t)
	db.Reset()

	queryMap := map[string]interface{}{
		"id": 0,
	}
	// Database is empty
	// Illegal Id 0
	content := createDummyContent(0)
	content.Question = "Test question"
	_, err := db.Add(content)
	if err != nil {
		t.Error(err)
	}
	k, err := db.Query(queryMap)
	if len(k) > 0 {
		t.Error("should not have found an instance with id 0")
	}
	k, err = db.Query(nil)
	queryMap["id"] = 1
	k, err = db.Query(queryMap)
	if len(k) != 1 {
		t.Error("Did not find instance with id 1")
	}

	// Add content again, should show no changes
	_, err = db.Add(content)
	if err != nil {
		t.Error(err)
	}
	k, err = db.Query(queryMap)
	if len(k) != 1 {
		t.Errorf("Expected still only one element in db, found %d", len(k))
	}


	db.Reset()
	err = db.Populate(5)
	if err != nil {
		t.Error(err.Error())
	}
	// add twice, should only be in there once after
	content.Id = 5
	_, err = db.Add(content)
	if err != nil {
		t.Error(err)
	}
	_, err = db.Add(content)
	if err != nil {
		t.Error(err)
	}
	queryMap["id"] = content.Id
	queryMap["lang"] = content.Language
	k, err = db.Query(queryMap)
	if err != nil {
		t.Error(err.Error())
	}
	if len(k) != 1 {
		t.Error("It should only be added once!")
	}

	// 10 skipped a few, should be 7
	content.Id = 10
	_, err = db.Add(content)
	if err != nil {
		t.Error(err)
	}
	queryMap["id"] = 6
	k, err = db.Query(queryMap)
	if len(k) != 1 {
		t.Errorf("Expected id to be 6")
	}

	// Close DB
	_=db.Close()

}


func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}
package main

import (
	"testing"
	"time"
)


func TestCreateDB(t *testing.T) {
	db := getDB()
	data, err := db.query(map[string]string{"lang": "all"})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 4 {
		t.Fatal("Could not get all instances, was not initialized correctly")
	}
}

func TestAddInstance(t *testing.T) {
	db := getDB()
	newInstance := Instance{
		Content{5, "Question #5?", "yeaaaaah."},
		Language{"en"},
		JSONTime{time.Now()},
	}
	numberOfInstances := db.GetLength()
	db.addInstance(newInstance)
	if numberOfInstances != db.GetLength()-1 {
		t.Error("Instance was not added to database")
	}
}

func TestRemoveById(t *testing.T) {
	db := getDB()
	var oldInstances = db.GetLength()

	// Remove valid
	db.removeById(1)
	if db.GetLength() == oldInstances {
		t.Error("did not remove instance")
	}

	oldInstances = db.GetLength()

	// Remove not existing should not change anything
	db.removeById(10)
	if db.GetLength() != oldInstances {
		t.Error("did remove instance but this one did not exist")
	}

}
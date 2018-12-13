package main

import (
	"testing"
	"time"
)


func TestAddInstance(t *testing.T) {
	db := getDB()
	newInstance := Instance{
		Content{5, "Question #5?", "yeaaaaah."},
		Languages.ES,
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
	db.removeById(1, Languages.EN)
	if db.GetLength() == oldInstances {
		t.Error("did not remove instance")
	}

	oldInstances = db.GetLength()

	// Remove not existing should not change anything
	db.removeById(10, "es")
	if db.GetLength() != oldInstances {
		t.Error("did remove instance but this one did not exist")
	}

}

func TestUpdateById(t *testing.T) {
	db := getDB()
	upgraded := Instance{
		Content{1, "How is life these days? Upgraded version", "So good"},
		Languages.EN,
		JSONTime{time.Now()},
	}
	db.updateById(1, "en", upgraded)
	d, _ := db.getById(1, Languages.EN)
	if d != upgraded {
		t.Error("The instance in the db was not updated")
	}

}

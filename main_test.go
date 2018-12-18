package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var db *DB

func populateDB(db *DB, instances int) {
	collection := db.client.Database(Database).Collection(Collection)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	for i := 0; i < instances; i++ {
		content := createDummyContent(i)
		_, err := collection.InsertOne(ctx, content)
		if err != nil {
			panic(err)
		}
	}
}

func setupDB(t *testing.T) {
	var err error
	if db == nil {
		Database = "testing"
		data := DB{}

		err = data.connect(MongoIp, MongoPort)
		if err != nil {
			t.Error(err)
		}
		data.reset()
		populateDB(&data, 10)
		if err != nil {
			t.Fatal(err)
		}
		db = &data
	}
}

func TestRootHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RootHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetHandler(t *testing.T) {
	tt := []struct {
		lang       string
		id         string
		shouldPass bool
	}{
		{"de", "1", true},
		{"en", "2", true},
		{"es", "-1", false},
		{"", "", true},
		{"de", "", true},
		{"", "1", true},
		{"asdf", "1", false},
	}

	setupDB(t)

	for _, tc := range tt {
		req, err := http.NewRequest("GET", "/content", nil)
		if err != nil {
			t.Error(err)
		}

		q := req.URL.Query()
		q.Add("lang", tc.lang)
		q.Add("id", tc.id)

		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router(db)
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on lang %s and id %s: got %v want %v",
				tc.lang, tc.id, rr.Code, http.StatusOK)
		}

	}
}

func TestDeleteHandler(t *testing.T) {
	setupDB(t)

	content := createDummyContent(1)
	// First add something that we will delete next

	conn := db.client.Database(Database).Collection(Collection)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	res, err := conn.InsertOne(ctx, content)
	if res.InsertedID == nil || err != nil {
		t.Error("could not put content in DB")
	}

	var other Content
	err = conn.FindOne(context.Background(), bson.M{"lang": content.Language, "id": content.Id}).Decode(&other)

	if err != nil {
		t.Error("Did not find instance")
	}

	for i := 0; i < 2; i++ {
		path := fmt.Sprintf("/content/%s/%d", content.Language, content.Id)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router(db)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("handler should have returned NoContent")
			print(rr.Code)
		}
	}
}

func TestPostHandler(t *testing.T) {
	setupDB(t)

	content := createDummyContent(123)
	js, err := json.Marshal(content)
	if err != nil {
		t.Error(err)
	}
	// First request
	req, err := http.NewRequest("POST", "/content", strings.NewReader(string(js)))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router := Router(db)
	router.ServeHTTP(rr, req)

	code := rr.Code
	if !(code == http.StatusNoContent || code == http.StatusCreated) {
		t.Error("handler did not pass on first POST")
	}

	// Now the instance is created => must return NoContent
	router.ServeHTTP(rr, req)
	if code = rr.Code; code != http.StatusNoContent {
		t.Errorf("handler did not pass on second POST of same instance,  returned %d ", code)
	}

}

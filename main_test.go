package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alex1ai/ugr-master-cc/authentication"
	. "github.com/alex1ai/ugr-master-cc/data"
	"github.com/mongodb/mongo-go-driver/bson"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

var db *DB

func populateDB(db *DB, instances int) {
	collection := db.Client.Database(Database).Collection(Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	for i := 0; i < instances; i++ {
		content := createDummyContent(i)
		_, err := collection.InsertOne(ctx, content)
		if err != nil {
			panic(err)
		}
	}
}

func setupDB(t *testing.T) {
	if db == nil {
		Database = "testing"
		data := DB{}
		portI, err := strconv.Atoi(MongoPort)
		err = data.Connect(MongoIp, portI)
		if err != nil {
			t.Error(err)
		}
		data.Reset()
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
	handler := http.HandlerFunc(StatusHandler)

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

	conn := db.Client.Database(Database).Collection(Collection)
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

func TestLoginHandler(t *testing.T) {
	setupDB(t)
	db.Reset()
	_, err := authentication.AddUserIfNotThere("test1234", "test1234", db)
	if err != nil {
		t.Error(err)
	}
	users := []struct {
		U    authentication.User
		Pass bool
	}{
		{authentication.User{Name: "abc", Password: "test1234"}, false},
		{authentication.User{Name: "test1234", Password: "test1234"}, true},
		{authentication.User{Name: "test1234", Password: "test124"}, false},
	}
	for _, user := range users {

		js, err := json.Marshal(user.U)
		// First request
		req, err := http.NewRequest("POST", "/login", strings.NewReader(string(js)))
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		router := Router(db)
		router.ServeHTTP(rr, req)

		tokenString := rr.Body.String()

		if len(tokenString) == 0 && user.Pass {
			t.Errorf("Valid user %s reveiced illegal tokenString, found %s", user.U, tokenString)
		} else if len(tokenString) > 0 && !user.Pass && rr.Code != http.StatusForbidden {
			t.Errorf("Unvalid user %s received legal tokenString, found %s with status %d", user.U, tokenString, rr.Code)
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/mongodb/mongo-go-driver/mongo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var client *mongo.Client

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

func setupDB(t *testing.T) {
	var err error
	client, _ = initializeDatabase(MongoIp, MongoPort)
	if err != nil {
		t.Fatal(err)
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
		{"asdf", "1", true},
	}
	setupDB(t)
	for _, tc := range tt {
		path := fmt.Sprintf("/content?lang=%s&id=%s", tc.lang, tc.id)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router(client)
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on lang %s and id %s: got %v want %v",
				tc.lang, tc.id, rr.Code, http.StatusOK)
		}

	}
}

func TestDeleteByIdHandler(t *testing.T) {
	for i := 0; i < 2; i++ {
		path := fmt.Sprintf("/content/%s/%d", "en", 1)
		req, err := http.NewRequest("DELETE", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := Router(client)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler should have return bad status on second time delete")
		}
		// TODO: Fix data management then, this will work. right now every requests creates a new DB
		//if rr.Code == http.StatusNotFound && i == 0 {
		//	t.Errorf("handler should have return OK on first time delete")
		//}
	}

}

func TestAddInstanceHandler(t *testing.T) {
	i := InstancePackage{
		{Content{4, "2 Nueva Edicion!", "Creo que no!"},
			Languages.ES,
			JSONTime{time.Now()},
		},
		{Content{5, "2 Nueva Edicion!", "Creo que no!"},
			Languages.ES,
			JSONTime{time.Now()},
		},
	}
	js, _ := json.Marshal(i)

	req, err := http.NewRequest("POST", "/content", strings.NewReader(string(js)))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := Router(client)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler did not pass on POST %s", "/content")
	}

	var resp JSONResponse
	dec := json.NewDecoder(rr.Body)

	if err = dec.Decode(&resp); err != nil {
		t.Error(err)
	}
	if cmp.Equal(i, resp.Data) {
		t.Errorf("Found different instances even though we updated an existing one\nFound %s\n expected %s", rr.Body.String(), js)
	}
}

func TestPostByIdHandler(t *testing.T) {
	i := InstancePackage{
		{Content{2, "2 Nueva Edicion!", "Creo que no!"},
			Languages.ES,
			JSONTime{time.Now()},
		},
	}
	js, _ := json.Marshal(i)

	req, err := http.NewRequest("POST", "/content", strings.NewReader(string(js)))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := Router(client)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler did not pass on POST %s", "/content")
	}

	var resp JSONResponse
	dec := json.NewDecoder(rr.Body)

	if err = dec.Decode(&resp); err != nil {
		t.Error(err)
	}
	if cmp.Equal(i, resp.Data) {
		t.Errorf("Found different instances even though we updated an existing one\nFound %s\n expected %s", rr.Body.String(), js)
	}

}

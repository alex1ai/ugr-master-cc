package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
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

	// Check the response body is what we expect.
	expected := `Hello There!`
	if strings.Trim(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetByIdHandler(t *testing.T) {
	tt := []struct {
		variable   string
		shouldPass bool
	}{
		{"1", true},
		{"2", true},
		{"3", true},
		{"4", false},
		{"asdf", false},

	}

	for _, tc := range tt {
		path := fmt.Sprintf("/get/%s", tc.variable)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Error(err)
		}

		rr := httptest.NewRecorder()

		// Need to create a router that we can pass the request through so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/get/{id}", GetByIdHandler)
		router.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on routeVariable %s: got %v want %v",
				tc.variable, rr.Code, http.StatusOK)
		}

	}
}

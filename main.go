package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Content struct {
	Question string
	Answer   string
}

func dummyContents() map[uint]Content {
	m := make(map[uint]Content)
	m[1] = Content{"How is life these days?", "So good"}
	m[2] = Content{"Are 2 questions sufficient?", "I do not think so!"}
	m[3] = Content{"Are 3 questions sufficient?", "I think so!"}
	return m
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello There!\n"))
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	j, _ := json.Marshal(dummyContents())
	w.Write(j)
}

func GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	content, ok := dummyContents()[uint(id)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "There is no id #%d", id)
	} else {
		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(content)
		w.Write(j)
	}
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/all", GetAllHandler).Methods("GET")
	r.HandleFunc("/get/{id:[0-9]+}", GetByIdHandler).Methods("GET")

	r.HandleFunc("/", RootHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

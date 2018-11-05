package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Type Modelling
type Content struct {
	Id       uint
	Question string
	Answer   string
}

type Language struct {
	Code string
}

type Instance struct {
	Content   Content
	Language  Language
	CreatedAt JSONTime
}

type JSONTime struct {
	time.Time
}

type JSONResponse struct {
	Status string
	Data InstancePackage
}

// Alias for Array
type InstancePackage []Instance

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%d\"", t.Unix())
	return []byte(stamp), nil
}


func (t *JSONTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	newTime, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*t = JSONTime{time.Unix(int64(newTime), 0)}
	return nil
}

// This will be later part of DB-API
func dummyInstances() InstancePackage {
	return InstancePackage{
		{
			Content{1, "How is life these days?", "So good"},
			Language{"en"},
			JSONTime{time.Now()},
		},
		{
			Content{2, "Are 2 questions sufficient?", "I do not think so!"},
			Language{"en"},
			JSONTime{time.Now()},
		},
		{
			Content{3, "Are 3 questions sufficient?", "I think so!"},
			Language{"en"},
			JSONTime{time.Now()},
		},
		{
			Content{2, "2 preguntas son suficiente?", "Creo que no!"},
			Language{"es"},
			JSONTime{time.Now()},
		},
	}
}

func filterByLanguage(langCode string) (ret InstancePackage) {
	for _, i := range dummyInstances() {
		if i.Language.Code == langCode {
			ret = append(ret, i)
		}
	}
	return ret
}

// END DB-API

// Helper functions
func jsonWrapper(status int, data InstancePackage) []byte {
	j, err := json.Marshal(JSONResponse{http.StatusText(status), data})
	if err != nil {
		return jsonWrapper(http.StatusBadRequest, InstancePackage{})
	}
	return j
}

// Route-Handlers

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the infogration REST-API!\n"))
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonWrapper(http.StatusOK, dummyInstances()))
}

func GetByLangHandler(w http.ResponseWriter, r *http.Request) {
	code, _ := mux.Vars(r)["code"]
	data := filterByLanguage(code)
	var status int
	if len(data) == 0 {
		status = http.StatusBadRequest
	} else {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	w.Write(jsonWrapper(status, data))
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/all", GetAllHandler).Methods("GET")
	r.HandleFunc("/get/{code}", GetByLangHandler).Methods("GET")
	r.HandleFunc("/", RootHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

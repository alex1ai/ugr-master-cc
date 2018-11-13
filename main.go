package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func getDB() DummyData {
	db := new(DummyData)
	err := db.create()
	if err != nil {
		log.Fatal(err)
	}
	return *db
}

// Helper functions
func jsonWrapper(status int, data InstancePackage) []byte {
	j, err := json.Marshal(JSONResponse{http.StatusText(status), data})
	if err != nil {
		return jsonWrapper(http.StatusBadRequest, InstancePackage{})
	}
	return j
}

// Route-Handlers

func sendResponse(writer http.ResponseWriter, status int, data InstancePackage) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(jsonWrapper(status, data))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, nil)
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	data, err := db.query(map[string]string{
		"lang": "all",
	})
	if err == nil {
		sendResponse(w, http.StatusOK, data)
	} else {
		sendResponse(w, http.StatusBadRequest, nil)
	}
}

func GetByLangHandler(w http.ResponseWriter, r *http.Request) {
	code, _ := mux.Vars(r)["code"]
	db := getDB()
	data, err := db.query(map[string]string{
		"lang": code,
	})
	var status int
	if len(data) == 0 || err != nil {
		status = http.StatusBadRequest
	} else {
		status = http.StatusOK
	}
	sendResponse(w, status, data)
}

func DeleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	db := getDB()
	idNumber, err := strconv.Atoi(id)
	if err != nil || idNumber < 0{

	}
	err = db.removeById(uint(idNumber))
	if err != nil {
		sendResponse(w, http.StatusNotFound, nil)
	}
	sendResponse(w, http.StatusOK, nil)
}

func Router() *mux.Router {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/all", GetAllHandler).Methods("GET")
	r.HandleFunc("/get/{code}", GetByLangHandler).Methods("GET")
	r.HandleFunc("/get/{id}", DeleteByIdHandler).Methods("DELETE")

	r.HandleFunc("/", RootHandler)

	return r
}

func main() {

	// Get Router
	r := Router()

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

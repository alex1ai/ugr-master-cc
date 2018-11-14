package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
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

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	sendResponse(w, http.StatusOK, nil)
}

func GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	db := getDB()
	data, err := db.getByLanguage("all")
	if err == nil {
		sendResponse(w, http.StatusOK, data)
	} else {
		sendResponse(w, http.StatusBadRequest, nil)
	}
}

func GetByLangHandler(w http.ResponseWriter, r *http.Request) {
	code, _ := mux.Vars(r)["lang"]
	db := getDB()
	data, err := db.getByLanguage(code)
	var status int
	if len(data) == 0 || err != nil {
		status = http.StatusBadRequest
	} else {
		status = http.StatusOK
	}
	sendResponse(w, status, data)
}

func PostByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	lang, _ := mux.Vars(r)["lang"]

	idInt, err := strconv.Atoi(id)
	db := getDB()

	data, err := db.getById(uint(idInt), lang)
	var status int
	if err == nil {
		status = http.StatusBadRequest
	} else {
		status = http.StatusOK
	}
	sendResponse(w, status, InstancePackage{data,})
}

func DeleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	lang, _ := mux.Vars(r)["lang"]
	db := getDB()
	idNumber, err := strconv.Atoi(id)
	if err != nil || idNumber < 0 {

	}
	err = db.removeById(uint(idNumber), lang)
	if err != nil {
		sendResponse(w, http.StatusNotFound, nil)
	}
	sendResponse(w, http.StatusOK, nil)
}

func AddInstanceHandler(w http.ResponseWriter, r *http.Request) {
	id, e1 := mux.Vars(r)["id"]
	lang, e2 := mux.Vars(r)["lang"]
	q, e3 := mux.Vars(r)["q"]
	a, e4 := mux.Vars(r)["a"]
	idI, err := strconv.Atoi(id)

	if !(e1 == e2 == e3 == e4) || err != nil || idI < 0 {
		sendResponse(w, http.StatusBadRequest, nil)
	}
	inst := Instance{
		Content{uint(idI), q, a},
		Language{lang},
		JSONTime{time.Now()},
	}
	db := getDB()

	db.addInstance(inst)
	sendResponse(w, http.StatusOK, InstancePackage{inst,})
}

func Router() *mux.Router {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/all", GetAllHandler).Methods("GET")
	r.HandleFunc("/content/{lang}", GetByLangHandler).Methods("GET")
	r.HandleFunc("/content/{lang}/{id}", DeleteByIdHandler).Methods("DELETE")
	r.HandleFunc("/content/{lang}/{id}", PostByIdHandler).Methods("POST")
	r.HandleFunc("/content/{lang}/{id}/{q}/{a}", AddInstanceHandler).Methods("PUT")

	r.HandleFunc("/", RootHandler)

	return r
}

func main() {

	// Get Router
	r := Router()

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

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
	db := getDB()
	decoder := json.NewDecoder(r.Body)
	status := http.StatusOK
	var ip InstancePackage
	if err := decoder.Decode(&ip); err != nil {
		status, ip = http.StatusBadRequest, InstancePackage{}
	} else {
		for _, j := range ip {
			id, lang := j.Content.Id, j.Language
			_, exists := db.getById(id, lang)
			// This (lang,id) is not yet known -> Put
			if !exists {
				db.addInstance(j)
			} else { // It is already known -> Update
				db.updateById(id, lang, j)
			}
		}
	}
	sendResponse(w, status, ip)
}

func DeleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	lang, _ := mux.Vars(r)["lang"]
	db := getDB()
	idNumber, err := strconv.Atoi(id)
	if err != nil || idNumber < 0 {
		sendResponse(w, http.StatusBadRequest, nil)
	}
	err = db.removeById(uint(idNumber), lang)
	if err != nil {
		sendResponse(w, http.StatusNotFound, nil)
	}
	sendResponse(w, http.StatusOK, nil)
}

func AddInstanceHandler(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	decoder := json.NewDecoder(r.Body)
	status := http.StatusOK
	var ip InstancePackage
	if err := decoder.Decode(&ip); err != nil {
		status, ip = http.StatusBadRequest, InstancePackage{}
	} else {
		for _, j := range ip {
			if _, exists := db.getById(j.Content.Id, j.Language); !exists {
				db.addInstance(j)
			}
		}
	}
	sendResponse(w, status, ip)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This is not the page you are looking for", http.StatusNotFound)
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetAllHandler).Methods("GET")
	r.HandleFunc("/content/{lang}", GetByLangHandler).Methods("GET")
	r.HandleFunc("/content/{lang}/{id}", DeleteByIdHandler).Methods("DELETE")
	r.HandleFunc("/content", PostByIdHandler).Methods("POST")
	r.HandleFunc("/content", AddInstanceHandler).Methods("PUT")
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/{.*}", ErrorHandler)

	return r
}

func main() {

	// Get Router
	r := Router()

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":443", r))
}

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

const (
	MongoPort = "27017"
	MongoIp   = "localhost"
	DEBUG     = true

	LangRegex = "^\\w{2}$"
	IdRegex   = "^[1-9][0-9]*"
)

var (
	Database   = "info"
	Collection = "content"
)

func Router(db *DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetHandler(db)).Methods("GET").
		Queries("lang", "{lang}", "id", "{id:[0-9]*}")
	r.HandleFunc("/content", GetHandler(db)).Methods("GET")

	r.HandleFunc("/content", PostPutHandler(db)).Methods("POST")
	r.HandleFunc("/content", PostPutHandler(db)).Methods("PUT")
	r.HandleFunc("/content/{lang}/{id}", DeleteHandler(db)).Methods("DELETE")

	r.HandleFunc("/init", InitHandler(db))
	r.HandleFunc("/reset", ResetHandler(db))
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	r.HandleFunc("/", StatusHandler)
	r.HandleFunc("/status", StatusHandler).Methods("GET")
	r.HandleFunc("/{.*}", ErrorHandler)

	return r
}

func main() {

	port := getEnv("PORT", "3000")
	ip := getEnv("IP", "0.0.0.0")
	home := os.Getenv("HOME")
	logPath := getEnv("LOG_FILE", home+"/logfile")
	logFile := initLogger(logPath)
	defer func() {
		_ = logFile.Close()
	}()

	// Initialize Datebase
	db := DB{}
	mPort := getEnv("MONGO_PORT", MongoPort)
	mPortI, err := strconv.Atoi(mPort)
	mIp := getEnv("MONGO_IP", MongoIp)
	err = db.connect(mIp, mPortI)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer func() {
		_ = db.close()
	}()

	r := Router(&db)

	// Add middleware
	r.Use(LoggingMiddleware)
	r.Use(LoggedInMiddleware)
	//loggedRouter := handlers.LoggingHandler(logFile, r)
	log.Infof("Starting web server on %s:%s", ip, port)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

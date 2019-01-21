package main

import (
	"fmt"
	"github.com/alex1ai/ugr-master-cc/authentication"
	"github.com/alex1ai/ugr-master-cc/data"
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

func Router(db *data.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetHandler(db)).Methods("GET").
		Queries("lang", "{lang}", "id", "{id:[0-9]*}")
	r.HandleFunc("/content", GetHandler(db)).Methods("GET")

	r.HandleFunc("/content", PostPutHandler(db)).Methods("POST")
	r.HandleFunc("/content", PostPutHandler(db)).Methods("PUT")
	r.HandleFunc("/content/{lang}/{id}", DeleteHandler(db)).Methods("DELETE")

	r.HandleFunc("/init", InitHandler(db))
	r.HandleFunc("/reset", ResetHandler(db))
	r.HandleFunc("/login", LoginHandler(db)).Methods("POST")

	r.HandleFunc("/", StatusHandler)
	r.HandleFunc("/status", StatusHandler).Methods("GET")
	r.HandleFunc("/{.*}", ErrorHandler)

	return r
}

func main() {

	port := getEnv("API_PORT", "3000")
	ip := getEnv("API_IP", "0.0.0.0")
	home := os.Getenv("HOME")
	logPath := getEnv("LOG_FILE", home+"/logfile")
	logFile := initLogger(logPath)
	defer func() {
		_ = logFile.Close()
	}()

	// Initialize Datebase
	db := data.DB{}
	mPort := getEnv("MONGO_PORT", MongoPort)
	mPortI, err := strconv.Atoi(mPort)
	if err != nil {
		log.Fatalf("Specified MONGO_PORT is not a number, found: %s", mPort)
		panic(err)
	}
	mIp := getEnv("MONGO_IP", MongoIp)
	err = db.Connect(mIp, mPortI)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	r := Router(&db)

	// Add middleware
	r.Use(LoggingMiddleware)
	r.Use(LoggedInMiddleware)
	log.Infof("Starting web server on %s:%s", ip, port)
	log.Infof("Using Mongo server on: %s:%d", mIp, mPortI)

	// Make sure admin user is registered
	_, err = authentication.RegisterAdmin(&db)

	if err != nil {
		log.Error("Could not register admin user or database problems")
	}

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

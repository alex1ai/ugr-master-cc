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
)

func checkAndSetEnvSecrets(){
	vars := map[string]string {
		"ADMIN_PW": "admin",
		"JWT_SECRET": "admin secret, change me!",
	}
	for k,v := range vars {
		set := os.Getenv(k)
		if set == "" {
			log.Warnf("You did not set env variable %s, using default value. Don't do this", k)
			if err := os.Setenv(k, v); err != nil {
				log.Fatalf("Could not set default env variable %s", k)
			}

		}
	}
}

func Router(db *data.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetHandler(db)).Methods("GET").
		Queries("lang", "{lang}", "id", "{id:[0-9]*}")
	r.HandleFunc("/content", GetHandler(db)).Methods("GET")

	r.HandleFunc("/content", EditContentHandler(db)).Methods("POST")
	r.HandleFunc("/content", EditContentHandler(db)).Methods("PUT")
	r.HandleFunc("/content/{lang}/{id}", DeleteHandler(db)).Methods("DELETE")

	r.HandleFunc("/init", InitHandler(db))
	r.HandleFunc("/reset", ResetHandler(db))
	r.HandleFunc("/login", LoginHandler(db)).Methods("POST")

	r.HandleFunc("/", StatusHandler)
	r.HandleFunc("/status", StatusHandler).Methods("GET")
	r.HandleFunc("/{.*}", ErrorHandler)

	// Add middleware
	r.Use(LoggingMiddleware)
	r.Use(LoggedInMiddleware)

	return r
}

func main() {
	// Specify where to run webservice
	port := getEnv("API_PORT", "3000")
	ip := getEnv("API_IP", "0.0.0.0")

	// Set Logfile
	home := os.Getenv("HOME")
	logPath := getEnv("LOG_FILE", home+"/logfile")
	logFile := initLogger(logPath)
	defer func() {
		_ = logFile.Close()
	}()

	// Check secret env variables or set defaults
	checkAndSetEnvSecrets()

	// Initialize Datebase
	db := data.DB{}
	defer func() {
		_ = db.Close()
	}()
	mPort := getEnv("MONGO_PORT", MongoPort)
	mPortI, err := strconv.Atoi(mPort)
	if err != nil {
		log.Fatalf("Specified MONGO_PORT is not a number, found: %s", mPort)
		os.Exit(1)
	}
	mIp := getEnv("MONGO_IP", MongoIp)
	err = db.Connect(mIp, mPortI)

	if err != nil {
		log.Fatalf("Error while connecting to DB: %s", err)
		os.Exit(1)
	}

	r := Router(&db)

	// Make sure admin user is registered
	_, err = authentication.RegisterAdmin(&db)
	if err != nil {
		log.Error("Could not register admin user or database problems")
	}


	log.Infof("Starting web server on %s:%s", ip, port)
	log.Infof("Using Mongo server on: %s:%d", mIp, mPortI)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

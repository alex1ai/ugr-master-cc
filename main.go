package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	MongoPort = 27017
	MongoIp   = "localhost"
	DEBUG     = true
	LangRegex = "^[a-z]{2}$"
	IdRegex   = "^[1-9][0-9]*"
)

var (
	Database   = "info"
	Collection = "content"
)

func initializeDatabase(ip string, port int) (client *mongo.Client, err error) {
	log.Infof("Connecting to Mongo Database, make sure it is running on %s:%d", ip, port)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%d", MongoIp, MongoPort))

	// Test reaching the DB
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	return
}

func sendResponse(writer http.ResponseWriter, status int, data []byte) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(data)
}

// ROUTES FOR WEBSERVICE
func RootHandler(w http.ResponseWriter, _ *http.Request) {
	sendResponse(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}

func GetHandler(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		id := r.FormValue("id")

		collection := c.Database(Database).Collection(Collection)

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		idOk, idEmpty := validateId(id)
		langOk, langEmpty := validateLang(lang)

		query := make(map[string]interface{})
		if idEmpty && langEmpty {
			query = nil
		}
		if (!idOk && !idEmpty) || (!langOk && !langEmpty) {
			http.Error(w, "Bad Parameters", http.StatusBadRequest)
		}
		if idOk {
			idi, _ := strconv.Atoi(id)
			query["id"] = uint(idi)
		}
		if langOk {
			query["lang"] = lang
		}

		cur, err := collection.Find(ctx, query)
		errorPanic(w, err)

		defer cur.Close(ctx)

		response := make([]Content, 1)
		for cur.Next(ctx) {
			var result Content
			err := cur.Decode(&result)
			errorPanic(w, err)
			log.Debug(result.Language)
			response = append(response, result)
		}
		errorPanic(w, cur.Err())
		j, err := json.Marshal(response)
		errorPanic(w, err)

		w.WriteHeader(http.StatusOK)
		w.Write(j)

	}
}

func PostPutHandler(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var instance Content
		if err := decoder.Decode(&instance); err != nil || !instance.validate() {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		id, lang := instance.Id, instance.Language

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		coll := c.Database(Database).Collection(Collection)

		opts := options.ReplaceOptions{}
		// If lang and id not there yet, this represents the same as a PUT request, i.e. creating a new Document
		opts.SetUpsert(true)
		_, err := coll.ReplaceOne(ctx, bson.M{"lang": lang, "id": id}, instance, &opts)
		errorPanic(w, err)

		sendResponse(w, http.StatusNoContent, nil)

	}
}

func DeleteHandler(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]
		lang, _ := mux.Vars(r)["lang"]
		idNumber, err := strconv.Atoi(id)
		if err != nil || idNumber < 0 {
			sendResponse(w, http.StatusBadRequest, nil)
		}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		coll := c.Database(Database).Collection(Collection)

		del, err := coll.DeleteOne(ctx, bson.M{"lang": lang, "id": uint(idNumber)})
		log.Debug(del.DeletedCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		sendResponse(w, http.StatusNoContent, nil)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This is not the page you are looking for", http.StatusNotFound)
	log.Warnf("Page not found: %s", r.URL.Path)
}

func Router(client *mongo.Client) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetHandler(client)).Methods("GET").
		Queries("lang", "{lang}", "id", "{id:[0-9]*}")
	r.HandleFunc("/content", GetHandler(client)).Methods("GET")

	r.HandleFunc("/content", PostPutHandler(client)).Methods("POST")
	r.HandleFunc("/content", PostPutHandler(client)).Methods("PUT")

	r.HandleFunc("/content/{lang}/{id}", DeleteHandler(client)).Methods("DELETE")

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/{.*}", ErrorHandler)

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RequestURI)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {

	port := getEnv("PORT", "3000")
	ip := getEnv("IP", "0.0.0.0")
	home := os.Getenv("HOME")
	logPath := getEnv("LOG_FILE", home+"/logfile")
	logFile := initLogger(logPath)
	defer logFile.Close()

	// Initialize Datebase
	client, err := initializeDatabase(MongoIp, MongoPort)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer client.Disconnect(context.Background())

	r := Router(client)

	// Add middleware
	r.Use(loggingMiddleware)

	log.Infof("Starting web server on %s:%s", ip, port)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

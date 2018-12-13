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
	"regexp"
	"strconv"
	"time"
)

const (
	MongoPort  = 27017
	MongoIp    = "localhost"
	DEBUG      = true

	LangRegex = "^[a-z]{2}$"
	IdRegex   = "^[1-9][0-9]*"
)

var (
	Database = "info"
	Collection = "content"
)

func initLogger(fileName string) *os.File {
	if !DEBUG {
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(file)
		return file
	}
	log.SetLevel(log.DebugLevel)
	f := os.File{}
	return &f
}

func initializeDatabase(ip string, port int) (client *mongo.Client, err error) {
	log.Infof("Connecting to Mongo Database, make sure it is running on %s:%d", ip, port)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%d", MongoIp, MongoPort))

	// Test reaching the DB
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	return
}

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
		j = jsonWrapper(http.StatusBadRequest, InstancePackage{})
	}
	log.WithFields(log.Fields{
		"status": status,
		"data":   string(j),
	}).Info("HTML Response")
	return j
}

// Route-Handlers

func sendResponse(writer http.ResponseWriter, status int, data []byte) {
	writer.Header().Set("ContentNew-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(data)
}

func errorPanic(w http.ResponseWriter, err error) {
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateId(id string) (matches bool, empty bool) {
	ok, err := regexp.MatchString(IdRegex, id)
	if err != nil {
		log.Debug(err.Error())
	}
	return ok, id == ""
}

func validateLang(lang string) (matches bool, empty bool) {
	ok, err := regexp.MatchString(LangRegex, lang)
	if err != nil {
		log.Debug(err.Error())
	}
	return ok, lang == ""
}

func populateDB(c *mongo.Client, instances int) {
	collection := c.Database(Database).Collection(Collection)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	for i := 0; i < instances; i++{
		content := createDummyContent()
		_, err := collection.InsertOne(ctx, content)
		if err != nil {
			log.Debug(err)
		}
	}

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

		response := make([]ContentNew, 1)
		for cur.Next(ctx) {
			var result ContentNew
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

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\"}"))
}

func PostPutHandler(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var instance ContentNew
		if err := decoder.Decode(&instance); err != nil || !instance.validate(){
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

func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func dropDB(client *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client.Database(Database).Drop(ctx)
	log.Debug("Dropped database")
}

func main() {

	port := getEnv("PORT", "3000")
	ip := getEnv("IP", "0.0.0.0")
	home := os.Getenv("HOME")
	logPath := getEnv("LOG_FILE", home+"/logfile.log")
	logFile := initLogger(logPath)
	defer logFile.Close()

	// Initialize Datebase
	client, err := initializeDatabase(MongoIp, MongoPort)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer client.Disconnect(context.Background())

	//dropDB(client)
	// Randomly init database TODO: delete
	//for i := 0; i < 10; i++{
	//	populateDB(client)
	//}
	//log.Debug("Created 10 instances")
	//os.Exit(0)
	// Get 	Router
	r := Router(client)

	// Add middleware
	r.Use(loggingMiddleware)

	log.Infof("Starting web server on %s:%s", ip, port)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

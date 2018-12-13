package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
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
	Database   = "info"
	Collection = "content"
	DEBUG      = true

	LangRegex = "^\\w{2}"
	IdRegex   = "^[1-9][0-9]*"
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

func sendResponse(writer http.ResponseWriter, status int, data InstancePackage) {
	writer.Header().Set("ContentNew-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(jsonWrapper(status, data))
}

func errorPanic(w http.ResponseWriter, err error) {
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ContentNew struct {
	Question  string `bson:"q, omitempty"`
	Answer    string `bson:"a, omitempty"`
	Id        uint   `bson:"id, omitempty"`
	Language  string `bson:"lang, omitempty"`
	Category  string `bson:"cat, omitempty"`
	CreatedAt uint   `bson:"time, omitempty"`
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

func populateDB(c *mongo.Client) {
	collection := c.Database(Database).Collection(Collection)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	content := ContentNew{"q1", "a1", 2, "en", "work", 1234}

	res, err := collection.InsertOne(ctx, content)
	log.Debug(res, err)

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
	sendResponse(w, http.StatusOK, nil)
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
	log.Warnf("Page not found: %s", r.URL.Path)
}

func Router(client *mongo.Client) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetHandler(client)).
		Methods("GET").
		Queries("lang", "{lang}", "id", "{id:[0-9]*}")
	r.HandleFunc("/content", GetHandler(client)).Methods("GET")
	r.HandleFunc("/content/{lang}/{id}", DeleteByIdHandler).Methods("DELETE")
	r.HandleFunc("/content", PostByIdHandler).Methods("POST")
	r.HandleFunc("/content", AddInstanceHandler).Methods("PUT")
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

	//populateDB(client)
	// Get 	Router
	r := Router(client)

	// Add middleware
	r.Use(loggingMiddleware)

	log.Infof("Starting web server on %s:%s", ip, port)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

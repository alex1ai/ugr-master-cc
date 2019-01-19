package main

import (
	"encoding/json"
	"fmt"
	"github.com/alex1ai/ugr-master-cc/authentication"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func sendResponse(writer http.ResponseWriter, status int, data []byte) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(data)

}

// ROUTES FOR WEBSERVICE
func RootHandler(w http.ResponseWriter, _ *http.Request) {
	sendResponse(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}

func GetHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		id := r.FormValue("id")

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

		response, err := db.query(query)

		j, err := json.Marshal(response)
		errorPanic(w, err)

		w.WriteHeader(http.StatusOK)
		w.Write(j)

	}
}

// TODO: If this is a put request, automatically fill Id-Number according to maximum in database
func PostPutHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var instance Content
		if err := decoder.Decode(&instance); err != nil || !instance.validate() {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		id, lang := instance.Id, instance.Language

		query := map[string]interface{}{
			"lang": lang,
			"id":   id,
		}

		_, err := db.update(query, instance)
		errorPanic(w, err)

		sendResponse(w, http.StatusNoContent, nil)

	}
}

func DeleteHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]
		lang, _ := mux.Vars(r)["lang"]
		idNumber, err := strconv.Atoi(id)
		if err != nil || idNumber < 0 {
			sendResponse(w, http.StatusBadRequest, nil)
		}

		query := map[string]interface{}{
			"lang": lang,
			"id":   uint(idNumber),
		}
		_, err = db.delete(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		sendResponse(w, http.StatusNoContent, nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user authentication.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tokenString, err := authentication.CreateToken(user.Name, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}
	_, err = w.Write([]byte(tokenString))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func initHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db.populate(10)
	}
}

func resetHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db.reset()
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This is not the page you are looking for", http.StatusNotFound)
	log.Warnf("Page not found: %s", r.URL.Path)
}

func Router(db *DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", GetHandler(db)).Methods("GET").
		Queries("lang", "{lang}", "id", "{id:[0-9]*}")
	r.HandleFunc("/content", GetHandler(db)).Methods("GET")

	r.HandleFunc("/content", PostPutHandler(db)).Methods("POST")
	r.HandleFunc("/content", PostPutHandler(db)).Methods("PUT")

	r.HandleFunc("/content/{lang}/{id}", DeleteHandler(db)).Methods("DELETE")

	r.HandleFunc("/init", initHandler(db))
	r.HandleFunc("/reset", resetHandler(db))
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/{.*}", ErrorHandler)

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func LoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Header)
		if r.URL.Path != "/login" && r.Method != http.MethodGet {
			token := r.Header.Get("Authorization")
			if token != "" {
				token = strings.Split(token, "\\w")[1]
				_, ok := authentication.ValidateToken(token)

				if ok {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "False credentials", http.StatusForbidden)
				}
			}
		}
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
	db := DB{}
	mPort := getEnv("MONGO_PORT", MongoPort)
	mPortI, err := strconv.Atoi(mPort)
	mIp := getEnv("MONGO_IP", MongoIp)
	err = db.connect(mIp, mPortI)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer db.close()

	r := Router(&db)

	// Add middleware
	r.Use(loggingMiddleware)
	r.Use(LoggedInMiddleware)
	//loggedRouter := handlers.LoggingHandler(logFile, r)
	log.Infof("Starting web server on %s:%s", ip, port)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), r))
}

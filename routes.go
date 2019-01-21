package main

import (
	"encoding/json"
	"errors"
	"github.com/alex1ai/ugr-master-cc/authentication"
	. "github.com/alex1ai/ugr-master-cc/data"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func sendResponse(writer http.ResponseWriter, data []byte, status ...int) {
	writer.Header().Set("Content-Type", "application/json")
	if len(status) > 0 {
		writer.WriteHeader(status[0])
	}
	_, err := writer.Write(data)
	if err != nil {
		http.Error(writer, "Could not send error", http.StatusInternalServerError)
	}
}

func sendError(writer http.ResponseWriter, statusCode int, err error) {
	http.Error(writer, err.Error(), statusCode)
}

// ROUTES FOR WEBSERVICE
func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	sendResponse(w, []byte("{\"status\": \"OK\"}"))
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
			sendError(w, http.StatusBadRequest, errors.New("bad parameters in query"))
		}
		if idOk {
			idi, _ := strconv.Atoi(id)
			query["id"] = uint(idi)
		}
		if langOk {
			query["lang"] = lang
		}

		response, err := db.Query(query)

		j, err := json.Marshal(response)
		errorPanic(w, err)

		sendResponse(w, j)
	}
}

// TODO: If this is a put request, automatically fill Id-Number according to maximum in database
func PostPutHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var instance Content
		if err := decoder.Decode(&instance); err != nil {
			sendError(w, http.StatusBadRequest, err)
		}

		id, lang := instance.Id, instance.Language

		query := map[string]interface{}{
			"lang": lang,
			"id":   id,
		}

		_, err := db.Update(query, instance)
		errorPanic(w, err)

		sendResponse(w, nil, http.StatusNoContent)

	}
}

func DeleteHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]
		lang, _ := mux.Vars(r)["lang"]
		idNumber, err := strconv.Atoi(id)
		if err != nil || idNumber < 0 {
			sendError(w, http.StatusBadRequest, err)
		}

		query := map[string]interface{}{
			"lang": lang,
			"id":   uint(idNumber),
		}
		_, err = db.Delete(query)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err)
		}
		sendResponse(w, nil, http.StatusNoContent)
	}
}

func LoginHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var user authentication.User
		if err := decoder.Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		tokenString, err := authentication.CreateToken(user.Name, user.Password, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		_, err = w.Write([]byte(tokenString))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func InitHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Populate(10)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err)
		}
	}
}

func ResetHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db.Reset()
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This is not the page you are looking for", http.StatusNotFound)
	log.Warnf("Page not found: %s", r.URL.Path)
}

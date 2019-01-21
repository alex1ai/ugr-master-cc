package main

import (
	"github.com/alex1ai/ugr-master-cc/authentication"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func LoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Header)
		defer func() {
			if err := recover(); err != nil {
				log.Fatal(err)
				http.Error(w, "Request was not correct", http.StatusBadRequest)
			}
		}()
		if r.URL.Path != "/login" && r.Method != http.MethodGet {
			token := r.Header.Get("Authorization")
			if token != "" {
				token = strings.TrimSpace(token)
				log.Debug(token)

				token = strings.Split(token, "\\w")[1]
				_, ok := authentication.ValidateToken(token)

				if ok {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "False credentials", http.StatusUnprocessableEntity)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

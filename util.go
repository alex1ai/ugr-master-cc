package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

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

func initLogger(fileName string) *os.File {
	t := time.Now().String()
	if !DEBUG {
		file, err := os.OpenFile(fmt.Sprintf("%s-%s.log", fileName, t), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func createDummyContent(id int) Content {
	langs := []string{"de", "en", "es", "ar"}
	lang := langs[rand.Intn(len(langs))]
	created := time.Now()
	return Content{"test 1", "test1 answer", uint(id), lang, "work", created}
}

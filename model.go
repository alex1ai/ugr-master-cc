package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Language Enum more or less
type LangRegistry struct {
	DE  string
	EN  string
	ES  string
	AR  string
	ALL string
}

type ContentNew struct {
	Question  string    `bson:"q, omitempty" json:"question"`
	Answer    string    `bson:"a, omitempty" json:"answer"`
	Id        uint      `bson:"id, omitempty" json:"id"`
	Language  string    `bson:"lang, omitempty" json:"lang"`
	Category  string    `bson:"cat, omitempty" json:"category"`
	CreatedAt time.Time `bson:"time, omitempty" json:"created_at"`
}

func createDummyContent() ContentNew{
	id := rand.Intn(20)
	langs := []string{"de", "en", "es", "ar"}
	lang := langs[rand.Intn(len(langs))]
	created := time.Now()
	return ContentNew{"test 1", "test1 answer", uint(id), lang, "work", created}
}

func (c *ContentNew) validate() bool {
	return true
}

func newLangRegistry() *LangRegistry {
	return &LangRegistry{
		DE:  "de",
		EN:  "en",
		ES:  "es",
		AR:  "ar",
		ALL: "all",
	}
}

var Languages = newLangRegistry()

// Type Modelling
type Content struct {
	Id       uint   `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Instance struct {
	Content   Content  `json:"content"`
	Language  string   `json:"language"`
	CreatedAt JSONTime `json:"createdAt"`
}

type JSONTime struct {
	time.Time
}

type JSONResponse struct {
	Status string          `json:"status"`
	Data   InstancePackage `json:"data"`
}

// Alias for Array
type InstancePackage []Instance

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%d\"", t.Unix())
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	newTime, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*t = JSONTime{time.Unix(int64(newTime), 0)}
	return nil
}

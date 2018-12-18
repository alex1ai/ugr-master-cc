package main

import (
	"math/rand"
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

type Content struct {
	Question  string    `bson:"q, omitempty" json:"question"`
	Answer    string    `bson:"a, omitempty" json:"answer"`
	Id        uint      `bson:"id, omitempty" json:"id"`
	Language  string    `bson:"lang, omitempty" json:"lang"`
	Category  string    `bson:"cat, omitempty" json:"category"`
	CreatedAt time.Time `bson:"time, omitempty" json:"created_at"`
}

func createDummyContent() Content {
	id := rand.Intn(20)
	langs := []string{"de", "en", "es", "ar"}
	lang := langs[rand.Intn(len(langs))]
	created := time.Now()
	return Content{"test 1", "test1 answer", uint(id), lang, "work", created}
}

func (c *Content) validate() bool {
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

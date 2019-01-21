package data

import (
	"time"
)

type Content struct {
	Question  string    `bson:"q, omitempty" json:"question, omitempty"`
	Answer    string    `bson:"a, omitempty" json:"answer, omitempty"`
	Id        uint      `bson:"id, omitempty" json:"id, omitempty"`
	Language  string    `bson:"lang, omitempty" json:"lang, omitempty"`
	Category  string    `bson:"cat, omitempty" json:"category, omitempty"`
	CreatedAt time.Time `bson:"time, omitempty" json:"created_at, omitempty"`
}

func (c *Content) validate() bool {
	return true
}

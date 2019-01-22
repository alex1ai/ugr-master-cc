package data

import (
	"time"
)

type Content struct {
	Question  string    `bson:"q" json:"question"`
	Answer    string    `bson:"a" json:"answer"`
	Id        uint      `bson:"id, omitempty" json:"id, omitempty"`
	Language  string    `bson:"lang" json:"lang"`
	Category  string    `bson:"cat" json:"category"`
	CreatedAt time.Time `bson:"created_at" json:"created_at, omitempty"`
}

func (c *Content) SetCreatedAtNow(){
	c.CreatedAt = time.Now()
}

func (c *Content) SetID(id uint) {
	c.Id = id
}


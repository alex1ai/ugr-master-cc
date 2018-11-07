package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Type Modelling
type Content struct {
	Id       uint `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Language struct {
	Code string `json:"code"`
}

type Instance struct {
	Content   Content `json:"content"`
	Language  Language `json:"language"`
	CreatedAt JSONTime `json:"createdAt"`
}

type JSONTime struct {
	time.Time
}

type JSONResponse struct {
	Status string `json:"status"`
	Data InstancePackage `json:"data"`
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

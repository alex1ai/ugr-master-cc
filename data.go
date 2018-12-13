package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

const(
	MONGOIP   ="mongoIp"
	MONGOPORT ="mongoPort"
)

type DBQuery struct{
	Lang string
	Id uint
}

type DB interface {
	GetSize() int
	setup(map[string]string) (interface{}, error)
	add(...Instance) error
	update(Instance) error
	remove(Instance) error
	get(...DBQuery) (InstancePackage, error)
}

type Mongo struct{}

func (db Mongo) setup(pars map[string]string) (client interface{}, err error) {
	port := pars[MONGOPORT]
	ip:= pars[MONGOIP]
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%s", ip, port))
	return
}

type DummyData struct {
	instances InstancePackage
}

func (dd *DummyData) GetLength() int {
	return len(dd.instances)
}

func (dd *DummyData) create() error {
	dd.instances = InstancePackage{
		{
			Content{1, "How is life these days?", "So good"},
			Languages.EN,
			JSONTime{time.Now()},
		},
		{
			Content{2, "Are 2 questions sufficient?", "I do not think so!"},
			Languages.EN,
			JSONTime{time.Now()},
		},
		{
			Content{3, "Are 3 questions sufficient?", "I think so!"},
			Languages.EN,
			JSONTime{time.Now()},
		},
		{
			Content{2, "2 preguntas son suficiente?", "Creo que no!"},
			Languages.ES,
			JSONTime{time.Now()},
		},
	}
	return nil
}

func findInData(dd DummyData, fn func(instance Instance) bool) (InstancePackage, error) {
	var ret InstancePackage
	for _, i := range dd.instances {
		if fn(i) {
			ret = append(ret, i)
		}
	}
	if len(ret) == 0 {
		return InstancePackage{}, errors.New("could not find instance")
	}
	return ret, nil
}

func (dd *DummyData) getByLanguage(langCode string) (resp InstancePackage, error error) {
	if langCode == "all" {
		return dd.instances, nil
	}
	return findInData(*dd, func(instance Instance) bool {
		return instance.Language == langCode
	})
}

func (dd *DummyData) getById(id uint, lang string) (resp Instance, found bool) {
	data, err := findInData(*dd, func(instance Instance) bool {
		return instance.Language == lang && instance.Content.Id == id
	})
	if len(data) != 1 || err != nil {
		return Instance{}, false
	}
	return data[0], true
}

func (dd *DummyData) close() error {
	dd.instances = nil
	return nil
}

func (dd *DummyData) addInstance(other Instance) error {
	data, _ := findInData(*dd, func(instance Instance) bool {
		return instance == other
	})
	if len(data) > 0 {
		return dd.updateById(other.Content.Id, other.Language, other)
	}

	dd.instances = append(dd.instances, other)
	return nil
}

func (dd *DummyData) removeById(id uint, lang string) error {
	numBefore := dd.GetLength()
	data, err := findInData(*dd, func(instance Instance) bool {
		return instance.Language != lang && instance.Content.Id != id
	})
	if numBefore == len(data) || err != nil {
		return errors.New("nothing was deleted")
	}

	dd.instances = data
	return nil
}

func (dd *DummyData) updateById(id uint, lang string, updateInstance Instance) error {
	var element = -1
	for e, i := range dd.instances {
		if i.Content.Id == id && i.Language == lang {
			element = e
		}
	}
	if element == -1 {
		return errors.New(fmt.Sprintf("could not find elemend %d", id))
	}

	dd.instances[element] = updateInstance

	return nil
}

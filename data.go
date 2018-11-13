package main

import (
	"errors"
	"time"
)

type DBConnector interface {
	create() error
	query(params map[string]string) (resp InstancePackage, error error)
	addInstance(newInstance Instance) error
	close() error
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
			Language{"en"},
			JSONTime{time.Now()},
		},
		{
			Content{2, "Are 2 questions sufficient?", "I do not think so!"},
			Language{"en"},
			JSONTime{time.Now()},
		},
		{
			Content{3, "Are 3 questions sufficient?", "I think so!"},
			Language{"en"},
			JSONTime{time.Now()},
		},
		{
			Content{2, "2 preguntas son suficiente?", "Creo que no!"},
			Language{"es"},
			JSONTime{time.Now()},
		},
	}
	return nil
}

func (dd *DummyData) query(params map[string]string) (resp InstancePackage, error error) {
	lang, ok := params["lang"]
	if !ok {
		return nil, errors.New("keyword lang not in params")
	}
	var ret InstancePackage
	if lang == "all" {
		return dd.instances, nil
	}
	for _, i := range dd.instances {
		if i.Language.Code == lang {
			ret = append(ret, i)
		}
	}

	return ret, nil
}

func (dd *DummyData) close() error {
	dd.instances = nil
	return nil
}

func (dd *DummyData) addInstance(other Instance) error {
	dd.instances = append(dd.instances, other)
	return nil
}

func (dd *DummyData) removeById(id uint) error{
	var ret InstancePackage
	numBefore := dd.GetLength()
	for _, i := range dd.instances {
		if i.Content.Id != id {
			ret = append(ret, i)
		}
	}
	if numBefore == len(ret) {
		return errors.New("nothing was deleted")
	}

	dd.instances = ret

	return nil
}

package main

import (
	"encoding/json"
	"errors"
	"sync"
)

type Definition struct {
	Origin      string
	Word        string
	Explanation string
}

func newDictionary(origin string) Dictionary {
	return Dictionary{Origin: origin, locker: &sync.Mutex{}}
}

type Dictionary struct {
	Origin      string
	Definitions []Definition
	locker      sync.Locker
}

func (d *Dictionary) add(word, explanation string) {
	d.locker.Lock()
	defer d.locker.Unlock()
	d.Definitions = append(d.Definitions, Definition{d.Origin, word, explanation})
}

func (d *Dictionary) addDefinition(definition Definition) {
	d.locker.Lock()
	defer d.locker.Unlock()
	d.Definitions = append(d.Definitions, definition)
}

func (d *Dictionary) toJson() ([]byte, error) {
	d.locker.Lock()
	defer d.locker.Unlock()
	return json.Marshal(d)
}

func (d *Dictionary) isEmpty() bool {
	return len(d.Definitions) == 0
}

func (d *Dictionary) size() int {
	return len(d.Definitions)
}

func (d *Dictionary) get(wordToLookFor string) (Definition, error) {
	d.locker.Lock()
	defer d.locker.Unlock()
	for _, definition := range d.Definitions {
		if definition.Word == wordToLookFor {
			return definition, nil
		}
	}
	return Definition{}, errors.New("Not found")
}

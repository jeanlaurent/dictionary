package main

import (
	"encoding/json"
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
	d.Definitions = append(d.Definitions, Definition{word, explanation, d.Origin})
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

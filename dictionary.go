package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

func initDictionary(origin string) Dictionary {
	dictionary := newDictionary(origin)
	loadWordFromFile(&dictionary)
	if dictionary.isEmpty() {
		loadSampleWords(&dictionary)
	}
	return dictionary
}

func loadWordFromFile(dictionary *Dictionary) {
	if _, err := os.Stat("dictionary.txt"); err == nil {
		data, err := ioutil.ReadFile("dictionary.txt")
		if err != nil {
			fmt.Println(err)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			splittedLine := strings.Split(line, ":")
			if len(splittedLine) == 2 {
				dictionary.add(splittedLine[0], splittedLine[1])
			}
		}
	}

}

func loadSampleWords(dictionary *Dictionary) {
	dictionary.add("foo", "is not bar")
	dictionary.add("bar", "is not foo")
	dictionary.add("qix", "is not foo nor bar")
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

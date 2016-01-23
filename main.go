package main

import (
	"encoding/json"
	"net/http"

	"fmt"
	"io/ioutil"

	"os"

	"strings"

	"github.com/gorilla/mux"
)

func searchWord(writer http.ResponseWriter, request *http.Request, dictionary *Dictionary) {
	query := request.URL.Query()["q"][0]
	definition, err := dictionary.get(query)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	payload, err := json.Marshal(definition)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Write(payload)
}

func putWord(writer http.ResponseWriter, request *http.Request, dictionary *Dictionary) {
	definition, err := readDefinition(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	dictionary.addDefinition(definition)
	writer.WriteHeader(http.StatusCreated)
}

func readDefinition(request *http.Request) (Definition, error) {
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return Definition{}, err
	}
	definition := &Definition{}
	err = json.Unmarshal(bytes, definition)
	if err != nil {
		return Definition{}, err
	}
	return *definition, nil
}

func listWords(writer http.ResponseWriter, request *http.Request, dictionary *Dictionary) {
	payload, err := dictionary.toJson()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Write(payload)
}

func withDictionary(dictionary *Dictionary, handler func(writer http.ResponseWriter, request *http.Request, dictionary *Dictionary)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handler(writer, request, dictionary)
	}
}

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	dictionary := newDictionary(hostname)
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
	if dictionary.isEmpty() {
		dictionary.add("foo", "is not bar")
		dictionary.add("bar", "is not foo")
		dictionary.add("qix", "is not foo nor bar")
	}
	fmt.Println("I'm",hostname, "I have a dictionary of", dictionary.size(), "words")

	router := mux.NewRouter()

	router.HandleFunc("/words", withDictionary(&dictionary, searchWord)).Queries("q", "{*}")
	router.HandleFunc("/words", withDictionary(&dictionary, listWords)).Methods("GET")
	router.HandleFunc("/words", withDictionary(&dictionary, putWord)).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./app/")))
	http.Handle("/", router)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

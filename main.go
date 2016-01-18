package main

import (
	"encoding/json"
	"net/http"

	"fmt"
	"io/ioutil"

	"github.com/gorilla/mux"
)

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
	router := mux.NewRouter()

	dictionary := newDictionary("itsme")
	dictionary.add("foo", "is not bar")
	dictionary.add("bar", "is not foo")
	dictionary.add("qix", "is not foo nor bar")

	router.HandleFunc("/", withDictionary(&dictionary, listWords)).Methods("GET")
	router.HandleFunc("/words", withDictionary(&dictionary, listWords)).Methods("GET")
	router.HandleFunc("/words", withDictionary(&dictionary, putWord)).Methods("POST")

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", router)
}

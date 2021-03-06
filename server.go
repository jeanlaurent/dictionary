package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const DEFAULT_PORT = 8080

func startServer(dictionary *Dictionary) {
	port := DEFAULT_PORT
	router := mux.NewRouter()

	router.HandleFunc("/words", withDictionary(dictionary, searchWord)).Queries("q", "{*}")
	router.HandleFunc("/words", withDictionary(dictionary, listWords)).Methods("GET")
	router.HandleFunc("/words", withDictionary(dictionary, putWord)).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./app/")))
	http.Handle("/", router)

	fmt.Println("Listening on", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

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

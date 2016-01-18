package main

import (
  "net/http"
  "encoding/json"
)

/*
var dictionary Dictionary

func listWords(writer http.ResponseWriter, request *http.Request) {
  payload, err := json.Marshal(dictionary)
  if  err != nil {
    http.Error(writer, err.Error() , http.StatusInternalServerError)
    return
  }
  writer.Write(payload)
}
*/

func main() {
  dictionary := Dictionary{Origin: "itsme"}
  dictionary.add("foo","is not bar")
  dictionary.add("bar","is not foo")
  dictionary.add("qix","is not foo nor bar")

  http.HandleFunc("/list", func (writer http.ResponseWriter, request *http.Request) {
    payload, err := json.Marshal(dictionary)
    if  err != nil {
      http.Error(writer, err.Error() , http.StatusInternalServerError)
      return
    }
    writer.Write(payload)
  })
  http.ListenAndServe(":8080", nil)
}

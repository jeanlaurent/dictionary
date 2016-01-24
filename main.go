package main

import (
	"fmt"

	"os"
)

func loadHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return hostname
}

func main() {

	hostname := loadHostName()

	dictionary := initDictionary(hostname)
	fmt.Println("I'm", hostname, "I have a dictionary of", dictionary.size(), "words")

	startServer(&dictionary)
}

package main

// import "shazam/shazam"

import (
	"log"
	"net/http"
	// "shazam/server"
)

// we need to be able to upgrade the connection from http to webSocket

func main() {

	var handler http.Handler
	http.Handle("/", handler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}

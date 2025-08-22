package main

// import "shazam/shazam"

import (
	"log"
	"net/http"
	"shazam/server"
)

// we need to be able to upgrade the connection from http to webSocket

func main() {
	http.HandleFunc("/", server.Handler)
	log.Println("starting server")

	log.Fatal(http.ListenAndServe(":1080", nil))

}

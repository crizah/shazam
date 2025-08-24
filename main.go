package main

// import "shazam/shazam"

import (
	"log"
	"net/http"
	"shazam/server"
)

// we need to be able to upgrade the connection from http to webSocket

func main() {

	// go func() {
	// 	for data := range server.PayloadChan {
	// 		log.Println("Received from client:", string(data))
	// 	}
	// }()
	http.HandleFunc("/", server.Handler)
	log.Println("starting server")

	// 1080 for getting bytes if soing

	log.Fatal(http.ListenAndServe(":1058", nil))

}

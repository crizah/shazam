package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SongInfo struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

var Songs = make(map[uint32]SongInfo) // resets everytime u run the code

func insertSongs(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight
	var songs map[uint32]SongInfo
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse body

	if err := json.NewDecoder(r.Body).Decode(&songs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Songs = songs
	for k, v := range songs {
		Songs[k] = v
	}

	fmt.Println(len(Songs))

	// fmt.Println("Received songs:", songs)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received successfully"))

	youtubeToMP3(Songs)
	fmt.Println("downloaded")
}

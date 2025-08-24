package server

// for getting song names

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type SongInfo struct {
// 	Name   string `json:"name"`
// 	Artist string `json:"artist"`
// 	Album  string `json:"album"`
// }

// var Songs = make(map[uint32]SongInfo) // resets everytime u run the code

// func getSongs(responseWriter http.ResponseWriter, req (*http.Request)) {
// 	// CORS headers

// 	webSocket, err := New(responseWriter, req)
// 	if err != nil {
// 		Errors = append(Errors, err)

// 	}

// 	webSocket.header.Set("Access-Control-Allow-Origin", "*")
// 	webSocket.header.Set("Access-Control-Allow-Methods", "POST, OPTIONS")

// 	webSocket.header.Set("Access-Control-Allow-Headers", "Content-Type")

// 	// Handle preflight
// 	var songs map[uint32]SongInfo
// 	if req.Method == http.MethodOptions {
// 		responseWriter.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	if req.Method != http.MethodPost {
// 		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse body

// 	if err := json.NewDecoder(req.Body).Decode(&songs); err != nil {
// 		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Songs = songs
// 	for k, v := range songs {
// 		Songs[k] = v
// 	}

// 	fmt.Println(len(Songs))

// 	// fmt.Println("Received songs:", songs)

// 	responseWriter.WriteHeader(http.StatusOK)
// 	responseWriter.Write([]byte("Received successfully"))

// }

// func insertSongs(w http.ResponseWriter, r *http.Request) {
// 	// CORS headers
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	// Handle preflight
// 	var songs map[uint32]SongInfo
// 	if r.Method == http.MethodOptions {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse body

// 	if err := json.NewDecoder(r.Body).Decode(&songs); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Songs = songs
// 	for k, v := range songs {
// 		Songs[k] = v
// 	}

// 	fmt.Println(len(Songs))

// 	// fmt.Println("Received songs:", songs)

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Received successfully"))

// }

// // http.HandleFunc("/get_songs", insertSongs)

// // 	fmt.Println("Server running at http://localhost:8080")
// // 	err := http.ListenAndServe(":8080", nil)
// // 	if err != nil {
// // 		fmt.Println("Error starting server:", err)
// // 	}

// // 	fmt.Println("fingerprint inserted into db")
// // 	// fmt.Printf(Songs[239205285].Name)

// // 	// youtubeToMP3(Songs)
// // 	// fmt.Printf("songs downloaded")

// // 	// Songs doesnt hod any value

// // 	fmt.Println(len(Songs))

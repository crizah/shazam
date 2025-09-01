package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"shazam/db"
	"shazam/server"
	"shazam/shazam"
)

var Errors []error
var ReceivedTracks []string
var ReceivedBytes []byte

func Handler(responseWriter http.ResponseWriter, req *http.Request) {

	webSocket, err := server.New(responseWriter, req)
	if err != nil {
		Errors = append(Errors, err)

	}

	err = webSocket.Handshake()
	if err != nil {
		Errors = append(Errors, err)
	}

	defer webSocket.Close()

	for {
		frame, err := webSocket.Recv()
		if err != nil {
			Errors = append(Errors, err)
			break
		}

		//   *  %x0 denotes a continuation frame

		//   *  %x1 denotes a text frame

		//   *  %x2 denotes a binary frame

		//   *  %x3-7 are reserved for further non-control frames

		//   *  %x8 denotes a connection close

		//   *  %x9 denotes a ping

		//   *  %xA denotes a pong

		//   *  %xB-F are reserved for further control frames

		switch frame.Opcode {

		case 0x1:

			fmt.Println("Received text:", string(frame.Payload))

			// unmarshal to get it into an array

			err := json.Unmarshal([]byte(string(frame.Payload)), &ReceivedTracks)
			if err != nil {
				Errors = append(Errors, err)
				return
			}

			Tracks := server.GetTracks(ReceivedTracks)
			SongIds := server.FindOnYoutube(Tracks) // find on youtube, download wav

			for i, sID := range SongIds {

				filename := fmt.Sprintf("%d.wav", sID)
				path := filepath.Join("C:\\Users\\shaiz\\Downloads\\shazam\\songs", filename)

				info, err := shazam.ReadWavFile(path)

				if err != nil {
					fmt.Print(err)
					return
				}

				PCMData := shazam.GetPCMData(info.Data)

				spectrogram := shazam.GetSpectrogram(PCMData, info.SampleRate, 5000.0, info.SampleRate/4)

				peaks := shazam.FindPeaks(spectrogram, info.AudioDuration)

				fp := shazam.GetFingerPrint(peaks, sID) // 442485 godamn fps
				fmt.Println(len(fp.Order))

				mongoClient, err := db.NewMongoClient()
				if err != nil {
					fmt.Println("client couldnt connect")
					fmt.Println(err)
					return
				}

				err = mongoClient.PutSongIds(sID, Tracks[i].Artist, Tracks[i].Name)

				if err != nil {

					fmt.Println("didnt register songIds ")
					fmt.Println(err)
					Errors = append(Errors, err)
					return

				}
				// working till here

				fmt.Println("check1 ")

				err = mongoClient.PutintoDB(fp) // error here

				if err != nil {
					// put a delete function here
					// delete all that have that sID

					fmt.Println("didnt register fps")

					fmt.Println(err)
					Errors = append(Errors, err)
					return
				}

				fmt.Println("added into db")

			}

		case 0x2: // binary frame
			fmt.Println("Received binary")
			ReceivedBytes = frame.Payload

			// processof raw audio data
			// find matches and send back to client

		case 0x8: // close
			fmt.Println("Client closed connection")
			return
		case 0x9: // ping
			fmt.Println("Received ping")

		case 0xA: // pong
			fmt.Println("Received pong")
		default:
			fmt.Println("Unknown opcode", frame.Opcode)
		}
	}

}

func main() {
	http.HandleFunc("/", Handler)
	log.Println("starting server")
	// 1080 for getting bytes of song
	// 1058 for getting playlist songs
	log.Fatal(http.ListenAndServe(":1058", nil))

}

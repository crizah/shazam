package shazam

// "fmt"
// "net/http"

// func main() {
// 	// s := []complex128{complex(1, 2), complex(3, 4), complex(5, 6)}
// 	// b := []complex128{complex(1, 2), complex(7, 9), complex(9, 10)}
// 	// var spectrogram [][]complex128
// 	// spectrogram = append(spectrogram, s)
// 	// spectrogram = append(spectrogram, b)
// 	// peaks := findPeaks(spectrogram, 0.5)
// 	// fmt.Println(len(peaks))
// 	// fmt.Println("peaks extracted")

// 	// fileName := "file_example_WAV_1MG.wav"
// 	// // header, err := extractHeader(fileName)
// 	// // if err != nil {
// 	// // 	fmt.Println(err)
// 	// // 	return

// 	// // }
// 	// // fmt.Println("header extracted")

// 	// PSM_data := getPSMData()

// 	// PCMsamples, err2 := readPSMdata(fileName, header)
// 	// if err2 != nil {
// 	// 	fmt.Println(err2)
// 	// 	return

// 	// }
// 	// fmt.Println(len(PCMsamples))
// 	// for i := 0; i < 10; i++ {
// 	// 	fmt.Print(PCMsamples[i], " ")
// 	// }
// 	// fmt.Println()

// 	// correct till here

// 	// spectrogram := getSpectrogram(PSM_data, int(header.SampleRate), 5000.0, int(header.SampleRate/4)) // error

// 	// fmt.Println("spectrogram made of size: ", len(spectrogram))

// 	// // err3 := saveSpectrogramAsPPM(spectrogram, "spectrogram_image.ppm")
// 	// // if err3 != nil {
// 	// // 	fmt.Println(err3)
// 	// // }

// 	// info := getMetaData(header)
// 	// peaks := findPeaks(spectrogram, info.audioDuration)
// 	// fmt.Println("number of peaks extracted: ", len(peaks))
// 	// fp := getFingerPrint(peaks, 12)
// 	// fmt.Println("fingerprint generation done", len(fp)) // should be same as peaks

// 	// // err4 := insertFPintoDB(fp)
// 	// if err4 != nil {
// 	// 	panic(err4)
// 	// }

// 	// fmt.Println("fingerprint inserted into db")
// 	// // kind of working
// 	// // assign songID based on idek what but yeah keep track of that
// 	// // song id = hash(song name || song artist || song album)

// 	// http.HandleFunc("/get_songs", insertSongs)

// 	// fmt.Println("Server running at http://localhost:8080")
// 	// err := http.ListenAndServe(":8080", nil)

// 	// if err != nil {
// 	// 	fmt.Println("Error starting server:", err)
// 	// }

// 	// fmt.Printf(Songs[239205285].Name)

// 	// youtubeToMP3(Songs)
// 	// fmt.Printf("songs downloaded")

// 	// Songs doesnt hod any value

// 	// fmt.Println(len(Songs))

// }

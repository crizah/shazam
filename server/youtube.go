package server

import (
	"fmt"
	"shazam/structs"
	"strings"

	"hash/fnv"
	"os/exec"
	"path/filepath"
	// "github.com/kkdai/youtube/v2"
)

func MakeSongID(track structs.Helper) uint32 {
	h := fnv.New32a()
	// Concatenate Name and Artist with a separator
	h.Write([]byte(track.Name + "|" + track.Artist))
	return h.Sum32()
}

// func MakeSongID(track string) uint32 {
// 	h := fnv.New32a()
// 	h.Write([]byte(track))
// 	return h.Sum32()
// }

func GetTracks(tracks []string) []structs.Helper {
	var result []structs.Helper
	// Received text: ["Comfortably Numb - Pink Floyd","Have a Cigar - Pink Floyd","Wearing the Inside Out - Pink Floyd"]

	for _, track := range tracks {
		// name and artist seperted by " - "

		parts := strings.SplitN(track, " - ", 2)

		var r structs.Helper

		r.Name = strings.TrimSpace(parts[0])
		r.Artist = strings.TrimSpace(parts[1])

		result = append(result, r)

	}

	return result

}

func FindOnYoutube(tracks []structs.Helper) []uint32 {

	var songIds []uint32
	for _, track := range tracks {

		sID := MakeSongID(track)

		filename := fmt.Sprintf("%d.wav", sID)
		outputPath := filepath.Join("C:\\Users\\shaiz\\Downloads\\shazam\\songs", filename) // or local ./downloads directory
		// chnage the poutput into tmp
		songIds = append(songIds, sID)

		// download into a cloud storage

		s := track.Name + " - " + track.Artist

		cmd := exec.Command("yt-dlp",
			"--extract-audio",
			"--audio-format", "wav",
			"--output", outputPath,
			"ytsearch1:"+s,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n%s\n", track, err, string(output))
			continue
		}

		fmt.Println("Downloaded:", outputPath)

	}
	fmt.Println("all songs downloaded")

	return songIds

}

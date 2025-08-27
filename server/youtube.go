package server

import (
	"fmt"

	"hash/fnv"
	"os/exec"
	"path/filepath"
	// "github.com/kkdai/youtube/v2"
)

func MakeSongID(track string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(track))
	return h.Sum32()

}

func FindOnYoutube(tracks []string) []uint32 {

	var ouputPaths []uint32
	for _, track := range tracks {

		sID := MakeSongID(track)

		filename := fmt.Sprintf("%d.wav", sID)
		outputPath := filepath.Join("C:\\Users\\shaiz\\Downloads\\shazam\\songs", filename) // or local ./downloads directory
		// chnage the poutput into tmp
		ouputPaths = append(ouputPaths, sID)

		// download into a cloud storage

		cmd := exec.Command("yt-dlp",
			"--extract-audio",
			"--audio-format", "wav",
			"--output", outputPath,
			"ytsearch1:"+track,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n%s\n", track, err, string(output))
			continue
		}

		fmt.Println("Downloaded:", outputPath)

	}
	fmt.Println("all songs downloaded")

	return ouputPaths

}

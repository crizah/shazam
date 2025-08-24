package server

import (
	"fmt"

	"os/exec"
	"path/filepath"
	// "github.com/kkdai/youtube/v2"
)

func FindOnYoutube(tracks []string) {
	for _, track := range tracks {

		filename := fmt.Sprintf("%s.mp3", track)
		outputPath := filepath.Join("C:\\Users\\shaiz\\Downloads\\shazam\\downloaded_wav", filename) // or local ./downloads directory

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

}

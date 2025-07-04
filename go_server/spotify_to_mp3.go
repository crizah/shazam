// songs -> youtube -> get

package main

import (
	"fmt"

	"os"
	"os/exec"
	"path/filepath"

	// "github.com/kkdai/youtube/v2"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func youtubeToMP3(songs map[uint32]SongInfo) {
	for id, song := range songs {
		query := song.Name + " - " + song.Artist
		fmt.Println("Searching YouTube for:", query)

		// Filepath where the MP3 will be saved
		filename := fmt.Sprintf("%d.mp3", id)
		outputPath := filepath.Join("/tmp", filename) // or local ./downloads directory

		// yt-dlp command
		cmd := exec.Command("yt-dlp",
			"--extract-audio",
			"--audio-format", "mp3",
			"--output", outputPath,
			"ytsearch1:"+query,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n%s\n", query, err, string(output))
			continue
		}

		fmt.Println("Downloaded:", outputPath)

		//S3
		err = uploadToS3(id, outputPath)
		if err != nil {
			fmt.Printf("Upload failed for %s: %v\n", filename, err)
		}
	}
}

func uploadToS3(songID uint32, filePath string) error {
	ctx := context.TODO()

	client, err := mongo.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	bucket, err := gridfs.NewBucket(
		client.Database("musicDB"),
	)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uploadStream, err := bucket.OpenUploadStream(fmt.Sprintf("%d.mp3", songID))
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write([]byte(filePath)) // or use io.Copy
	return err
}

// func convertMP3toWAV() {
// 	// for every mp3 in s3, convert it into WAV and save it, delete the MP3 version

// }

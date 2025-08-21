package shazam

// do error handelinmg properly
// figure out the songID generation thing
// clean up the Matched zStruct
// clean up all structs
// priority queue implementation

import (
	"math"
	"shazam/db"
	"sort"
)

// do the processing of the file
// finger printing on that process
// matching

// To perform a search, the above fingerprinting step is
// performed on a captured sample sound file to generate a set
// of hash:time offset records. Each hash from the sample is
// used to search in the database for matching hashes. For
// each matching hash found in the database, the
// corresponding offset times from the beginning of the
// sample and database files are associated into time pairs.
// The time pairs are distributed into bins according to the
// track ID associated with the matching database hash.
// After all sample hashes have been used to search in the
// database to form matching time pairs, the bins are scanned
// for matches. Within each bin the set of time pairs
// represents a scatterplot of association between the sample
// and database sound files. If the files match, matching
// features should occur at similar relative offsets from the
// beginning of the file, i.e. a sequence of hashes in one file
// should also occur in the matching file with the same
// relative time sequence. The problem of deciding whether a
// match has been found reduces to detecting a significant
// cluster of points forming a diagonal line within the
// scatterplot.

// NEED TO DO ALL ERROR HANDELING

type Candidate struct {
	SongID uint32
	Points int
}

func FindMatches(sample []float64, sampleRate int, audioDuration float64, songId uint32) (*[]Candidate, error) {

	spectrogram := getSpectrogram(sample, sampleRate, 5000.0, sampleRate/4)
	peaks := findPeaks(spectrogram, audioDuration)

	fp := GetFingerPrint(peaks, songId) // set of ALL hashes of the song sample
	// search fp in the database
	Bins, err := db.SearchDB(fp)

	if err != nil {
		return nil, err

	}

	var candidates []Candidate

	for id, matches := range Bins {

		// 		a sequence of hashes in one file
		// should also occur in the matching file with the same
		// relative time sequence

		// per song
		bin := make(map[uint32]db.Matched)

		for _, match := range matches {
			bin[match.MatchedHash] = match

		}

		var binOrder []uint32
		for _, fpKey := range fp.Order {
			_, ok := bin[fpKey]
			if ok {
				binOrder = append(binOrder, fpKey)
			}

		}

		// binOrder and fp.Order now have relative position of the values

		points := 0

		for i := 0; i < len(binOrder)-1; i++ {
			for j := i + 1; j < len(binOrder); j++ {

				diff1 := math.Abs(float64(fp.Order[i] - fp.Order[j]))
				diff2 := math.Abs(float64(binOrder[i] - binOrder[j]))
				if math.Abs(diff1-diff2) < 100 {
					points++
				}

			}

		}

		can := Candidate{SongID: id, Points: points}

		candidates = append(candidates, can)

	}

	// sort candidates based on their points in decreasing order

	sort.Slice(candidates[:], func(i, j int) bool {
		return candidates[i].Points > candidates[j].Points
	})

	return &candidates, nil

}

// For each Matched in bins[songId]:

// Compute offset := MatchedTime - SampleTime.

// Count how many times each offset (or offset bucket) occurs.

// Because times won’t be exactly equal, bucket them (e.g., round to nearest 0.1s).

// The largest bucket count = score for that candidate song.

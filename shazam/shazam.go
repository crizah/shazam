package shazam

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

func FindMatches(sample []float64, sampleRate int, audioDuration float64, songId uint32) error {

	spectrogram := getSpectrogram(sample, sampleRate, 5000.0, sampleRate/4)
	peaks := findPeaks(spectrogram, audioDuration)

	fp := GetFingerPrint(peaks, songId) // set of ALL hashes of the song sample
	// search fp in the database

	for h, info := range fp {
		matched := SearchDB(h, info.anchor_time)
	}

	// m := make(map[uint32]uint32)
	// for h, info := range fp {
	// 	m[h] = info.anchor_time
	// }

	// search in db

	return nil

}

type Matched struct {
	SampleTime  uint32
	MatchedTime uint32
	DBsongId    uint32
}

func SearchDB(hash uint32, sampleTime uint32) []Matched {

	// For
	// each matching hash found in the database, the
	// corresponding offset times from the beginning of the
	// sample and database files are associated into time pairs.
	// The time pairs are distributed into bins according to the
	// track ID associated with the matching database hash.

	var result Matched
	return result

}

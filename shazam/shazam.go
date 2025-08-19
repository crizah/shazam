package shazam

// do the processing of the file
// finger printing on that process
// matching

func FindMatches(sample []float64, sampleRate int, audioDuration float64, songId uint32) error {

	spectrogram := getSpectrogram(sample, sampleRate, 5000.0, sampleRate/4)
	peaks := findPeaks(spectrogram, audioDuration)

	fp := GetFingerPrint(peaks, songId)
	// search fp in the database
	m := make(map[uint32]uint32)
	for h, info := range fp {
		m[h] = info.anchor_time
	}

	return nil

}

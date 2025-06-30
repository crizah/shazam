package main

import "math/cmplx"

var band_ranges = []Band{
	{0, 10}, {10, 20}, {20, 40}, {40, 80}, {80, 160}, {160, 511}, // so weird i have to add a comma at the end here
}

func findPeaks(spectrogram [][]complex128, audioDuration float64) []Peak {
	// error here

	var peaks []Peak
	frameDuration := audioDuration / float64(len(spectrogram))

	for i, frame := range spectrogram {
		// per frame
		// 6 strongPoints per frame
		strongPoints_of_frame_i := make([]StrongPoint, len(band_ranges))
		for _, band := range band_ranges {
			var strongP StrongPoint
			// per band, one strongpoint
			var maxMag float64 = -82

			for j, freq := range frame[band.min:band.max] { // slices are so great. so helpful
				mag := cmplx.Abs(freq)

				if mag > maxMag {
					maxMag = mag
					realIdx := band.min + j
					strongP.freq = freq
					strongP.mag = mag
					strongP.freq_idx = realIdx

				}

			}

			strongPoints_of_frame_i = append(strongPoints_of_frame_i, strongP)
		}

		// get average and only add all greater than average
		sum := 0.0
		for _, s := range strongPoints_of_frame_i {
			sum += s.mag
		}

		avg := sum / float64(len(strongPoints_of_frame_i))

		// dont add the ones lesser than avg
		for _, s := range strongPoints_of_frame_i { // _ decalare and not use
			if s.mag >= avg {
				var peak Peak
				peak.Frequency = s.freq

				// calculate time
				a := float64(s.freq_idx) * frameDuration / float64(len(frame))
				b := float64(i) * frameDuration
				peak.Time = a + b
				peaks = append(peaks, peak)
			}
		}

	}

	return peaks

}

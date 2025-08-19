package shazam

import (
	"fmt"
	"math"

	"math/cmplx"
)

// pcm samples to specyrogram to peaks to fingerprints
//  downsize + fft

const (
	frameSize = 1024
	hopSize   = frameSize / 32
)

func LowPassFilter(PCM []float64, sampleRate int, cutoffFrequency float64) []float64 {
	// filetered := make([]float64, len(PCM)) // this was error, donot do this, idek how
	var filetered []float64
	rc := 1.0 / (2.0 * math.Pi * cutoffFrequency)
	dt := 1.0 / float64(sampleRate)
	alpha := dt / (rc + dt)

	// fmt.Println(rc, ", ", dt, ", ", alpha)

	prev := float64(PCM[0])

	for i := 0; i < len(PCM); i++ {
		a := alpha*PCM[i] + (1-alpha)*prev
		filetered = append(filetered, a)
		prev = a

	}

	return filetered
}

func downSample(filtered []float64, originalSR int, targetSR int) []float64 {

	var downSampled []float64

	ratio := originalSR / targetSR

	for i := 0; i < len(filtered); i += ratio {
		end := i + ratio
		if end > len(filtered) {
			end = len(filtered)
		}

		sum := float64(0)

		for j := i; j < end; j++ {
			sum += filtered[j]

		}

		avg := sum / float64(end-i)
		downSampled = append(downSampled, avg)

	}

	return downSampled

}

func hann() []float64 {
	window := make([]float64, frameSize)
	for i := range window {
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(frameSize-1)))
	}

	return window

}
func frameSignal(samples []float64) [][]complex128 {
	var frames [][]complex128
	window := hann()
	// fmt.Println(len(window))

	numFrames := len(samples) / (frameSize - hopSize)
	for i := 0; i < numFrames; i++ {
		start := i * hopSize
		end := start + frameSize

		if end > len(samples) {
			end = len(samples)
		}

		frame := make([]float64, frameSize)

		for j := 0; j < len(window); j++ { // size of window is frameSize
			frame[j] = samples[start+j] * window[j]
		}

		freq := make([]complex128, len(frame))

		for j := range frame {
			freq[j] = complex(frame[j], 0)
		}

		fft(freq)
		frames = append(frames, freq)
	}

	return frames

}

func getSpectrogram(PCMsamples []float64, sampleRate int, cutOffFreq float64, targetRate int) [][]complex128 {
	lps := LowPassFilter(PCMsamples, sampleRate, cutOffFreq)
	fmt.Println("lps size: ", len(lps))
	// fine till here

	downsampled := downSample(lps, sampleRate, targetRate)
	fmt.Println("downsampled size: ", len(downsampled))

	// fine till here

	spectrogram := frameSignal(downsampled) // error here

	return spectrogram
}

func Finalspectrogram() {

	spectrogram := getSpectrogram()
	peaks := findPeaks(spectrogram)
	fingerPrints := getFingerPrint(peaks)

}

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

const (
	r = 5
)

func compressHash(h Hash) uint32 {
	a := uint32(h.a_frequency<<23) | uint32(h.t_frequency<<14) | uint32(h.time)
	return a
}

func getFingerPrint(peaks []Peak, songId uint32) map[uint32]information {
	var fingerPrint = make(map[uint32]information)
	for i, anchor := range peaks {
		// per anchor
		for j := i + 1; j < i+r && j < len(peaks); j++ {
			// per a, t pair
			target := peaks[j]
			anchor_freq := int(real(anchor.Frequency))
			target_freq := int(real(target.Frequency))
			time_diff := uint32((target.Time - anchor.Time) * 1000)

			h := Hash{anchor_freq, target_freq, time_diff}
			// compress h to uint 32

			hash_i := compressHash(h)

			anchor_time := uint32(anchor.Time * 1000)
			info := information{anchor_time, songId}
			// intermediate := []uint32{anchor_time, songId}
			fingerPrint[hash_i] = info
		}
	}
	return fingerPrint
	// i think type error

}

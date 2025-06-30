package main

import (
	"fmt"
	"math"
)

// pcm samples to spectrogram
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

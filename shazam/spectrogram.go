package shazam

import (
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

	// https://dsp.stackexchange.com/questions/9425/how-to-determine-alpha-smoothing-constant-of-a-lpf

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
	// avg out the 4 and takje those

	var downSampled []float64

	ratio := originalSR / targetSR // 4

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

	// https://www.sciencedirect.com/topics/engineering/hanning-window
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

func GetSpectrogram(PCMsamples []float64, sampleRate int, cutOffFreq float64, targetRate int) [][]complex128 {
	lps := LowPassFilter(PCMsamples, sampleRate, cutOffFreq) // cutoff freq = 5000
	// fmt.Println("lps size: ", len(lps))
	// fine till here

	downsampled := downSample(lps, sampleRate, targetRate) // targetRate = original rate /4
	// fmt.Println("downsampled size: ", len(downsampled))

	// fine till here

	spectrogram := frameSignal(downsampled) // error here

	return spectrogram
	// rows freq, columns time and values are intensity

	// 	The translationinvariant aspect means that fingerprint hashes derived from
	// corresponding matching content are reproducible
	// independent of position within an audio file, as long as the
	// temporal locality containing the data from which the hash
	// is computed is contained within the file.
}

var band_ranges = []Band{
	{0, 10}, {10, 20}, {20, 50}, {50, 80}, {80, 150}, {150, 350}, {350, 520},
}

type Peak struct {
	Time      float64
	Frequency complex128
}

type StrongPoint struct {
	freq     complex128
	mag      float64
	freq_idx int
}

type Band struct {
	min, max int
}

func FindPeaks(spectrogram [][]complex128, audioDuration float64) []Peak {
	// error here

	// 	A time-frequency point is a candidate peak if it has a
	// higher energy content than all its neighbors in a region
	// centered around the point. Candidate peaks are chosen
	// according to a density criterion in order to assure that the
	// time-frequency strip for the audio file has reasonably
	// uniform coverage. The peaks in each time-frequency
	// locality are also chosen according amplitude, with the
	// justification that the highest amplitude peaks are most
	// likely to survive the distortions listed above.

	// max amplitudes in each locality

	// 20 - 5000hz are the freq
	// 440hz is what music is tuned to

	// var peaks []Peak
	// for i, frame := range spectrogram{

	// 	// per frame, 6 candidate points
	// 	var Sp Peak
	// 	for _, band := range band_ranges{
	// 		mm := -897986
	// 		var candidatePeaks StrongPoint

	// 		for j, freq := range frame{

	// 		}

	// 	}

	// rows freq, columns time and values are intensity

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

			for j, freq := range frame[band.min:band.max] {
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

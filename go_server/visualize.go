package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
)

// Normalize the spectrogram to uint8 values using log scale
func normalizeSpectrogram(spec [][]complex128) [][]uint8 {
	height := len(spec)
	width := len(spec[0])

	// Find min/max log magnitude
	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64

	logMag := make([][]float64, height)
	for i := range spec {
		logMag[i] = make([]float64, width)
		for j := range spec[i] {
			mag := math.Log10(1 + cmplx.Abs(spec[i][j]))
			logMag[i][j] = mag
			if mag < minVal {
				minVal = mag
			}
			if mag > maxVal {
				maxVal = mag
			}
		}
	}

	rangeVal := maxVal - minVal
	if rangeVal == 0 {
		rangeVal = 1
	}

	// Normalize to [0, 255]
	norm := make([][]uint8, height)
	for i := range logMag {
		norm[i] = make([]uint8, width)
		for j := range logMag[i] {
			norm[i][j] = uint8(255 * (logMag[i][j] - minVal) / rangeVal)
		}
	}

	return norm
}

// Clamp helper
func clamp(x, minVal, maxVal float64) float64 {
	if x < minVal {
		return minVal
	}
	if x > maxVal {
		return maxVal
	}
	return x
}

// Jet-style colormap mapping
func jetColorMap(value uint8) (uint8, uint8, uint8) {
	x := float64(value) / 255.0

	r := uint8(255 * clamp(math.Min(4*(x-0.75), 1.0), 0.0, 1.0))
	g := uint8(255 * clamp(math.Min(4*math.Abs(x-0.5)-1.0, 1.0), 0.0, 1.0))
	b := uint8(255 * clamp(math.Min(4*(0.25-x), 1.0), 0.0, 1.0))

	return r, g, b
}

// Save PPM (P3 format)
func saveSpectrogramAsPPM(spectrogram [][]complex128, filename string) error {
	normSpec := normalizeSpectrogram(spectrogram)

	height := len(normSpec)
	width := len(normSpec[0])

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write PPM header
	fmt.Fprintf(file, "P3\n%d %d\n255\n", width, height)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r, g, b := jetColorMap(normSpec[i][j])
			fmt.Fprintf(file, "%d %d %d  ", r, g, b)
		}
		fmt.Fprintln(file)
	}

	return nil
}

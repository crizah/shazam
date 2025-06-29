package main

import (
	"fmt"
)

func main() {
	s := []complex128{complex(1, 2), complex(3, 4), complex(5, 6)}
	b := []complex128{complex(1, 2), complex(7, 9), complex(9, 10)}
	var spectrogram [][]complex128
	spectrogram = append(spectrogram, s)
	spectrogram = append(spectrogram, b)
	peaks := findPeaks(spectrogram, 0.5)
	fmt.Println(len(peaks))
	fmt.Println("peaks extracted")

}

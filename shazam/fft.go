package shazam

import (
	"math"
)

func fft(a []complex128) {

	// f(t) = A cos (Bt + C)

	// https://cp-algorithms.com/algebra/fft.html

	n := len(a)

	if n <= 1 {
		return
	}

	// var even[n/2] complex128 // not working for some reason?
	even := make([]complex128, n/2)
	odd := make([]complex128, n/2)
	// var odd[n/2] complex128

	for i := 0; i < n/2; i++ {
		even[i] = a[2*i]
		odd[i] = a[2*i+1]
	}

	fft(even)
	fft(odd)

	for k := 0; k < n/2; k++ {

		t := complex((math.Cos(-2*math.Pi*float64(k)/float64(n))), (math.Sin(-2*math.Pi*float64(k)/float64(n)))) * odd[k]
		a[k] = even[k] + t
		a[k+n/2] = even[k] - t
	}

}

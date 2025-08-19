package shazam

type Information struct {
	anchor_time uint32
	songID      uint32
}

const (
	r = 5
)

func GetFingerPrint(peaks []Peak, songId uint32) map[uint32]Information {
	var fingerPrint = make(map[uint32]Information)
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
			info := Information{anchor_time, songId}
			// intermediate := []uint32{anchor_time, songId}
			fingerPrint[hash_i] = info
		}
	}
	return fingerPrint
	// save this to db

	// i think type error

}

func compressHash(h Hash) uint32 {
	a := uint32(h.a_frequency<<23) | uint32(h.t_frequency<<14) | uint32(h.time)
	return a
}

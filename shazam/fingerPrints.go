package shazam

import "shazam/structs"

const (
	r = 10
)

func GetFingerPrint(peaks []Peak, songId uint32) structs.OMap {

	// 	Fingerprint hashes are formed from the constellation map,
	// in which pairs of time-frequency points are combinatorially
	// associated. Anchor points are chosen, each anchor point
	// having a target zone associated with it. Each anchor point
	// is sequentially paired with points within its target zone,
	// each pair yielding two frequency components plus the time
	// difference between the points (Figure 1C and 1D). These
	// hashes are quite reproducible, even in the presence of noise
	// and voice codec compression. Furthermore, each hash can
	// be packed into a 32-bit unsigned integer. Each hash is also
	// associated with the time offset from the beginning of the
	// respective file to its anchor point, though the absolute time
	// is not a part of the hash itself.

	var fingerPrint structs.OMap
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

			info := structs.Information{anchor_time, songId}
			// intermediate := []uint32{anchor_time, songId}

			fingerPrint.Map[hash_i] = info
			fingerPrint.Order = append(fingerPrint.Order, hash_i)

		}
	}
	return fingerPrint
	// save this to db

}

type Hash struct {
	a_frequency int
	t_frequency int
	time        uint32
}

func compressHash(h Hash) uint32 {
	a := uint32(h.a_frequency<<23) | uint32(h.t_frequency<<14) | uint32(h.time)
	return a
}

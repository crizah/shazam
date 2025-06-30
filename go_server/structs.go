package main

type Hash struct {
	a_frequency int
	t_frequency int
	time        uint32
}

type information struct {
	anchor_time uint32
	songID      uint32
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

type WAVHeader struct {

	// https://docs.fileformat.com/audio/wav/

	Riff            [4]byte
	File_size       uint32
	Wave            [4]byte
	Fmt             [4]byte
	Length          uint32
	Format          uint16
	Channels        uint16
	SampleRate      uint32
	Byte_rate       uint32
	Block_allign    uint16
	Bits_per_sample uint16
	Data            [4]byte
	Data_size       uint32
}

type Info struct {
	sampleRate    int
	audioDuration float64
}

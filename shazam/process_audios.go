package main

// return PCM samples from WAV file. later nned to expand to convert any audio file to WAV format
// later tho

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type WAVHeader struct {

	// https://docs.fileformat.com/audio/wav/

	Riff            [4]byte
	File_size       uint32
	Wave            [4]byte
	Fmt             [4]byte
	Length          uint32
	Format          uint16
	Channels        uint16 // want this
	SampleRate      uint32 // want this
	Byte_rate       uint32 // wnat this
	Block_allign    uint16
	Bits_per_sample uint16 // want this
	Data            [4]byte
	Data_size       uint32 // wnt tihsu
}

func WriteWavHeader(f *os.File, data []byte, sampleRate int, channels int, bitsPerSample int) error {

	l_size := uint32(16)
	bytesPerSample := bitsPerSample / 8
	blockAlign := uint16(channels * bytesPerSample)
	subchunk2Size := uint32(len(data))

	header := WAVHeader{
		Riff:            [4]byte{'R', 'I', 'F', 'F'},
		File_size:       uint32(36 + len(data)),
		Wave:            [4]byte{'W', 'A', 'V', 'E'},
		Fmt:             [4]byte{'f', 'm', 't', ' '},
		Length:          l_size,
		Format:          uint16(1), // PCM format
		Channels:        uint16(channels),
		SampleRate:      uint32(sampleRate),
		Byte_rate:       uint32(sampleRate * channels * bytesPerSample),
		Block_allign:    blockAlign,
		Bits_per_sample: uint16(bitsPerSample),
		Data:            [4]byte{'d', 'a', 't', 'a'},
		Data_size:       subchunk2Size,
	}

	err := binary.Write(f, binary.LittleEndian, header) // write the header file
	return err

}

func PutHeaderIntoFile(filename string, data []byte, sampleRate int, channels int, bitsPerSample int) error {
	// var header WAVHeader

	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	err = WriteWavHeader(file, data, sampleRate, channels, bitsPerSample)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err

}

type Info struct {
	SampleRate int
	Channals   int
	Data       []byte

	AudioDuration float64
}

func ReadWavFile(filename string) (*Info, error) {

	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// header file needs to be 44 bytes
	if len(f) < 44 {
		return nil, errors.New("Wav file too small")
	}

	var header WAVHeader

	err = binary.Read(bytes.NewReader(f[:44]), binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	audioduration := float64(len(f)) / float64(int(header.Channels)*2*int(header.SampleRate))

	info := &Info{
		SampleRate:    int(header.SampleRate),
		Channals:      int(header.Channels),
		Data:          f[:44],
		AudioDuration: audioduration,
	}

	return info, nil

}

func getPSMData(input []byte) []float64 {

	n := int(len(input) / 2)
	output := make([]float64, n)

	for i := 0; i < n; i++ {
		start := i * 2
		sample := int16(binary.LittleEndian.Uint16(input[start : start+2]))
		output[i] = float64(sample)
	}

	// for i := 0; i < len(input); i += 2 {
	// 	// Interpret bytes as a 16-bit signed integer (little-endian)
	// 	sample := int16(binary.LittleEndian.Uint16(input[i : i+2]))

	// 	// Scale the sample to the range [-1, 1]
	// 	output[i/2] = float64(sample) / 32768.0
	// }

	return output

}

// func readPSMdata(filename string, header WAVHeader) ([]float64, error) {

// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open file: %v", err)
// 	}
// 	defer file.Close()

// 	// Seek to start of audio data
// 	_, err = file.Seek(int64(binary.Size(header)), 0)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to seek: %v", err)
// 	}

// 	bytesPerSample := header.Bits_per_sample / 8
// 	numSamples := int(header.Data_size) / int(bytesPerSample)

// 	// Read raw PCM bytes
// 	raw := make([]byte, header.Data_size)
// 	_, err = file.Read(raw)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read audio data: %v", err)
// 	}

// 	samples := make([]float64, numSamples)

// 	for i := 0; i < numSamples; i++ {
// 		start := i * 2
// 		sample := int16(binary.LittleEndian.Uint16(raw[start : start+2]))
// 		samples[i] = float64(sample)
// 	}

// 	return samples, nil
// }
